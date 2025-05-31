import { useState } from 'react';
import { formatDistanceToNow } from 'date-fns';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Comment } from '@/types';
import { useAuth } from '@/contexts/AuthContext';
import { LoginModal } from '@/components/LoginModal';
import RichTextEditor from '@/components/RichTextEditor';
import RichTextDisplay from '@/components/RichTextDisplay';
import { Loader2, MessageCircle } from 'lucide-react';

interface CommentCardProps {
  comment: Comment;
  onEdit?: (commentId: string, newContent: string) => void;
  onReply?: (parentId: string, content: string) => void;
  onLoadReplies?: (commentId: string) => Promise<Comment[]>;
  isReply?: boolean;
}

const CommentCard = ({ 
  comment, 
  onEdit, 
  onReply, 
  onLoadReplies, 
  isReply = false 
}: CommentCardProps) => {
  const { user, isAuthenticated } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const [isReplying, setIsReplying] = useState(false);
  const [editContent, setEditContent] = useState(comment.content);
  const [replyContent, setReplyContent] = useState('');
  const [showReplies, setShowReplies] = useState(false);
  const [showLoginModal, setShowLoginModal] = useState(false);
  const [loadingReplies, setLoadingReplies] = useState(false);

  const formatDate = (dateString: string) => {
    try {
      return formatDistanceToNow(new Date(dateString), { addSuffix: true });
    } catch {
      return 'Unknown time';
    }
  };

  const isOwner = user && comment.author?.id === user.id;

  const handleEdit = async () => {
    if (editContent.trim() && onEdit) {
      try {
        await onEdit(comment.id, editContent.trim());
        setIsEditing(false);
      } catch (error) {
        // Error is handled by the API service
      }
    }
  };

  const handleReply = async () => {
    if (replyContent.trim() && onReply) {
      try {
        await onReply(comment.id, replyContent.trim());
        setReplyContent('');
        setIsReplying(false);
      } catch (error) {
        // Error is handled by the API service
      }
    }
  };

  const handleLoadReplies = async () => {
    if (onLoadReplies && !loadingReplies) {
      setLoadingReplies(true);
      try {
        await onLoadReplies(comment.id);
        setShowReplies(true);
      } catch (error) {
        console.error('Error loading replies:', error);
      } finally {
        setLoadingReplies(false);
      }
    }
  };

  const handleEditClick = () => {
    if (!isAuthenticated) {
      setShowLoginModal(true);
      return;
    }
    if (!isOwner) {
      return; // Don't show edit for non-owners
    }
    setIsEditing(true);
  };

  const handleReplyClick = () => {
    if (!isAuthenticated) {
      setShowLoginModal(true);
      return;
    }
    setIsReplying(!isReplying);
  };

  const handleLoginSuccess = () => {
    setShowLoginModal(false);
    // The user can now proceed with their intended action
  };

  return (
    <>
      <div className={`${isReply ? 'ml-8 border-l-2 border-gray-200 pl-4' : ''}`}>
        <Card className="w-full border border-gray-200 hover:shadow-md transition-shadow">
          <CardContent className="p-4">
            <div className="flex justify-between items-start mb-3">
              <div className="flex items-center gap-3">
                <div className="flex-shrink-0">
                  {comment.author?.avatar_url ? (
                    <img
                      src={comment.author.avatar_url}
                      alt={comment.author.display_name || comment.author.username || 'User'}
                      className="w-8 h-8 rounded-full"
                    />
                  ) : (
                    <div className="w-8 h-8 rounded-full bg-purple-100 flex items-center justify-center">
                      <span className="text-purple-600 font-medium text-sm">
                        {comment.author?.display_name?.[0] || comment.author?.username?.[0] || "Anonymous"}
                      </span>
                    </div>
                  )}
                </div>
                
                <div className="flex flex-col">
                  <div className="flex items-center gap-2">
                    <span className="font-medium text-gray-900">
                      {comment.author?.display_name || comment.author?.username || 'Anonymous'}
                    </span>
                    <span className="text-xs text-gray-500">
                      {formatDate(comment.created_at)}
                    </span>
                  </div>
                </div>
              </div>
            </div>
            {isEditing ? (
              <div className="space-y-3">
                <RichTextEditor
                  value={editContent}
                  onChange={setEditContent}
                  placeholder="Edit your comment..."
                  className="min-h-[120px]"
                />
                <div className="flex gap-2">
                  <Button onClick={handleEdit} size="sm">
                    Save
                  </Button>
                  <Button 
                    onClick={() => {
                      setIsEditing(false);
                      setEditContent(comment.content);
                    }} 
                    variant="outline" 
                    size="sm"
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            ) : (
              <div>
                <div className="mb-3">
                  <RichTextDisplay content={comment.content} className="text-gray-700" />
                </div>
                
                <div className="flex items-center gap-2 text-sm">
                  {isOwner && (
                    <Button
                      onClick={handleEditClick}
                      variant="ghost"
                      size="sm"
                      className="h-8 px-2 text-gray-600 hover:text-blue-600"
                    >
                      <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                      </svg>
                      Edit
                    </Button>
                  )}
                  
                  <Button
                    onClick={handleReplyClick}
                    variant="ghost"
                    size="sm"
                    className="h-8 px-2 text-gray-600 hover:text-green-600"
                  >
                    <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
                    </svg>
                    Reply
                  </Button>
                  
                  {!showReplies && comment.replies_count > 0 && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={handleLoadReplies}
                      disabled={loadingReplies}
                      className="mt-2"
                    >
                      {loadingReplies ? (
                        <>
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                          Loading...
                        </>
                      ) : (
                        <>
                          <MessageCircle className="mr-2 h-4 w-4" />
                          Load Replies
                          <Badge variant="secondary" className="ml-2">
                            {comment.replies_count}
                          </Badge>
                        </>
                      )}
                    </Button>
                  )}

                  {showReplies && comment.children && comment.children.length > 0 && (
                    <Button
                      onClick={() => setShowReplies(false)}
                      variant="ghost"
                      size="sm"
                      className="h-8 px-2 text-gray-600 hover:text-purple-600"
                    >
                      <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
                      </svg>
                      Hide Replies ({comment.children.length})
                    </Button>
                  )}
                </div>
              </div>
            )}

            {isReplying && (
              <div className="mt-4 space-y-3 border-t border-gray-200 pt-4">
                <RichTextEditor
                  value={replyContent}
                  onChange={setReplyContent}
                  placeholder="Write a reply..."
                  className="min-h-[120px]"
                />
                <div className="flex gap-2">
                  <Button onClick={handleReply} size="sm" disabled={!replyContent.trim()}>
                    Reply
                  </Button>
                  <Button 
                    onClick={() => {
                      setIsReplying(false);
                      setReplyContent('');
                    }} 
                    variant="outline" 
                    size="sm"
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Replies Section */}
        {showReplies && comment.children && comment.children.length > 0 && (
          <div className="mt-4 space-y-3">
            <div className="flex items-center justify-between">
              <h4 className="text-sm font-medium text-gray-700">
                Replies ({comment.children.length})
              </h4>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setShowReplies(false)}
                className="text-gray-500 hover:text-gray-700"
              >
                Hide Replies
              </Button>
            </div>
            <div className="space-y-3 pl-4 border-l-2 border-gray-100">
              {comment.children.map((reply: Comment) => (
                <CommentCard
                  key={reply.id}
                  comment={reply}
                  onEdit={onEdit}
                  onReply={onReply}
                  onLoadReplies={onLoadReplies}
                  isReply={true}
                />
              ))}
            </div>
          </div>
        )}
      </div>

      <LoginModal
        isOpen={showLoginModal}
        onClose={() => setShowLoginModal(false)}
        onSuccess={handleLoginSuccess}
      />
    </>
  );
};

export { CommentCard }; 