import { useState, useEffect, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { PostsAPI, CommentsAPI } from "@/services/api";
import { Post, Comment, PaginationMeta } from "@/types";
import { useInfiniteScroll } from "@/hooks/useInfiniteScroll";
import { useAuth } from "@/contexts/AuthContext";
import { LoginModal } from "@/components/LoginModal";
import { toast } from "sonner";
import LoadingState from "@/components/LoadingState";
import ErrorState from "@/components/ErrorState";
import { Button } from "@/components/ui/button";
import PostHeader from "@/components/PostHeader";
import CommentForm from "@/components/CommentForm";
import CommentsList from "@/components/CommentList";

const PostComments = () => {
  const { postId } = useParams<{ postId: string }>();
  const navigate = useNavigate();
  const { user, isAuthenticated } = useAuth();

  // State management
  const [post, setPost] = useState<Post | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingComments, setLoadingComments] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [meta, setMeta] = useState<PaginationMeta | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isCommenting, setIsCommenting] = useState(false);
  const [newComment, setNewComment] = useState("");
  const [showLoginModal, setShowLoginModal] = useState(false);
  const [isEditingPost, setIsEditingPost] = useState(false);
  const [editPostTitle, setEditPostTitle] = useState("");
  const [editPostContent, setEditPostContent] = useState("");

  // API calls
  const loadPost = useCallback(async () => {
    if (!postId) return;

    try {
      setLoading(true);
      setError(null);
      const postData = await PostsAPI.getPost(postId);
      setPost(postData);
    } catch (err: any) {
      setError(err.message || "Failed to load post");
    } finally {
      setLoading(false);
    }
  }, [postId]);

  const loadComments = useCallback(
    async (page: number = 1, append: boolean = false) => {
      if (!postId) return;

      try {
        if (page === 1) {
          setLoadingComments(true);
        } else {
          setLoadingMore(true);
        }

        const response = await CommentsAPI.getComments(postId, page, 10);

        if (append) {
          setComments((prev) => [...prev, ...response.comments]);
        } else {
          setComments(response.comments);
        }

        setMeta(response.meta);
      } catch (err: any) {
        setError(err.message || "Failed to load comments");
      } finally {
        setLoadingComments(false);
        setLoadingMore(false);
      }
    },
    [postId]
  );

  const loadMoreComments = useCallback(() => {
    if (meta && meta.hasNext && !loadingMore) {
      loadComments(meta.page + 1, true);
    }
  }, [meta, loadingMore, loadComments]);

  // Infinite scroll hook
  useInfiniteScroll({
    hasNextPage: meta?.hasNext || false,
    isLoading: loadingMore,
    loadMore: loadMoreComments,
    threshold: 200,
  });

  // Event handlers
  const handleEditComment = async (commentId: string, newContent: string) => {
    try {
      const updatedComment = await CommentsAPI.updateComment(
        commentId,
        newContent
      );
      
      // Recursive function to update comments at any nesting level
      const updateCommentsWithEdit = (comments: Comment[]): Comment[] => {
        return comments.map((comment) => {
          if (comment.id === commentId) {
            return {
              ...comment,
              content: updatedComment.content,
              updated_at: updatedComment.updated_at,
            };
          } else if (comment.children && comment.children.length > 0) {
            return {
              ...comment,
              children: updateCommentsWithEdit(comment.children),
            };
          }
          return comment;
        });
      };

      setComments((prev) => updateCommentsWithEdit(prev));
    } catch (error) {
      // Error is handled by the API service
    }
  };

  const handleReply = async (parentId: string, content: string) => {
    if (!postId) return;

    try {
      const newReply = await CommentsAPI.createComment(
        postId,
        content,
        parentId
      );

      const updateCommentsWithReply = (comments: Comment[]): Comment[] => {
        return comments.map((comment) => {
          if (comment.id === parentId) {
            const updatedChildren = comment.children
              ? [newReply, ...comment.children]
              : [newReply];
            return {
              ...comment,
              replies_count: (comment.replies_count || 0) + 1,
              children: updatedChildren,
            };
          } else if (comment.children && comment.children.length > 0) {
            return {
              ...comment,
              children: updateCommentsWithReply(comment.children),
            };
          }
          return comment;
        });
      };

      setComments((prev) => updateCommentsWithReply(prev));
      toast.success("Reply posted successfully!");
    } catch (error) {
      console.error("Error posting reply:", error);
      toast.error("Failed to post reply");
    }
  };

  const handleLoadReplies = async (commentId: string): Promise<Comment[]> => {
    try {
      const repliesResponse = await CommentsAPI.getReplies(commentId, 1, 10);
      
      // Recursive function to update comments at any nesting level
      const updateCommentsWithReplies = (comments: Comment[]): Comment[] => {
        return comments.map((comment) => {
          if (comment.id === commentId) {
            return { ...comment, children: repliesResponse.comments };
          } else if (comment.children && comment.children.length > 0) {
            return {
              ...comment,
              children: updateCommentsWithReplies(comment.children),
            };
          }
          return comment;
        });
      };

      setComments((prev) => updateCommentsWithReplies(prev));
      return repliesResponse.comments;
    } catch (error) {
      throw error;
    }
  };

  const handleNewComment = async () => {
    if (!postId || !newComment.trim()) return;

    if (!isAuthenticated) {
      setShowLoginModal(true);
      return;
    }

    try {
      const comment = await CommentsAPI.createComment(
        postId,
        newComment.trim()
      );
      setComments((prev) => [comment, ...prev]);
      setNewComment("");
      setIsCommenting(false);
      toast.success("Comment posted successfully!");

      if (post) {
        setPost((prev) =>
          prev
            ? { ...prev, commentsCount: (prev.commentsCount || 0) + 1 }
            : null
        );
      }
    } catch (error) {
      console.error("Error posting comment:", error);
      toast.error("Failed to post comment");
    }
  };

  const handleCommentClick = () => {
    if (!isAuthenticated) {
      setShowLoginModal(true);
      return;
    }
    setIsCommenting(!isCommenting);
  };

  const handleEditPost = async () => {
    if (!post || !editPostTitle.trim() || !editPostContent.trim()) return;

    try {
      const updatedPost = await PostsAPI.updatePost(post.id, {
        title: editPostTitle.trim(),
        content: editPostContent.trim(),
      });
      setPost(updatedPost);
      setIsEditingPost(false);
      toast.success("Post updated successfully!");
    } catch (error) {
      console.error("Error updating post:", error);
      toast.error("Failed to update post");
    }
  };

  const handleEditPostClick = () => {
    if (!post) return;
    setEditPostTitle(post.title);
    setEditPostContent(post.content);
    setIsEditingPost(true);
  };

  const handleCancelPostEdit = () => {
    setIsEditingPost(false);
    setEditPostTitle("");
    setEditPostContent("");
  };

  const handleCancelComment = () => {
    setIsCommenting(false);
    setNewComment("");
  };

  const handleLoginSuccess = () => {
    setShowLoginModal(false);
  };

  const handleNavigateBack = () => {
    navigate("/dashboard");
  };

  // Computed values
  const isPostOwner = user && post && post.created_by === user.id;

  // Effects
  useEffect(() => {
    if (postId) {
      loadPost();
      loadComments(1, false);
    }
  }, [postId, loadPost, loadComments]);

  // Render loading state
  if (loading) {
    return <LoadingState />;
  }

  // Render error state
  if (error && !post) {
    return (
      <ErrorState
        error={error}
        onNavigateBack={handleNavigateBack}
        onRetry={loadPost}
      />
    );
  }

  return (
    <>
      <div className="min-h-screen bg-gray-50 py-8">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <Button
            onClick={handleNavigateBack}
            variant="ghost"
            className="flex items-center gap-2"
          >
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15 19l-7-7 7-7"
              />
            </svg>
            Back to Dashboard
          </Button>
          <div className="mb-6"></div>

          {post && (
            <PostHeader
              post={post}
              isPostOwner={isPostOwner!}
              isEditingPost={isEditingPost}
              editPostTitle={editPostTitle}
              editPostContent={editPostContent}
              onEditPostClick={handleEditPostClick}
              onEditPostTitle={setEditPostTitle}
              onEditPostContent={setEditPostContent}
              onSavePost={handleEditPost}
              onCancelEdit={handleCancelPostEdit}
            />
          )}

          <CommentForm
            isCommenting={isCommenting}
            newComment={newComment}
            meta={meta}
            onCommentClick={handleCommentClick}
            onNewCommentChange={setNewComment}
            onSubmitComment={handleNewComment}
            onCancelComment={handleCancelComment}
          />

          <CommentsList
            comments={comments}
            loadingComments={loadingComments}
            loadingMore={loadingMore}
            meta={meta}
            onEditComment={handleEditComment}
            onReply={handleReply}
            onLoadReplies={handleLoadReplies}
            onCommentClick={handleCommentClick}
          />
        </div>
      </div>

      <LoginModal
        isOpen={showLoginModal}
        onClose={() => setShowLoginModal(false)}
        onSuccess={handleLoginSuccess}
      />
    </>
  );
};

export default PostComments;
