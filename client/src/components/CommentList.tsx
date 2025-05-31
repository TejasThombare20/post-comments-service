import React from 'react';
import { Button } from '@/components/ui/button';
import { CommentCard } from '@/components/CommentCard';
import { Comment, PaginationMeta } from '@/types';
import CommentSkeleton from './CommentSkeleton';

interface CommentsListProps {
  comments: Comment[];
  loadingComments: boolean;
  loadingMore: boolean;
  meta: PaginationMeta | null;
  onEditComment: (commentId: string, newContent: string) => Promise<void>;
  onReply: (parentId: string, content: string) => Promise<void>;
  onLoadReplies: (commentId: string) => Promise<Comment[]>;
  onCommentClick: () => void;
}

const CommentsList: React.FC<CommentsListProps> = ({
  comments,
  loadingComments,
  loadingMore,
  meta,
  onEditComment,
  onReply,
  onLoadReplies,
  onCommentClick,
}) => {
  if (loadingComments && comments.length === 0) {
    return (
      <div className="space-y-6">
        {Array.from({ length: 5 }).map((_, index) => (
          <CommentSkeleton key={index} />
        ))}
      </div>
    );
  }

  if (comments.length === 0 && !loadingComments) {
    return (
      <div className="text-center py-12">
        <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg className="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
        </div>
        <h3 className="text-lg font-medium text-gray-900 mb-2">No comments yet</h3>
        <p className="text-gray-600 mb-4">Be the first to comment on this post!</p>
        <Button onClick={onCommentClick} variant="outline">
          Add the first comment
        </Button>
      </div>
    );
  }

  return (
    <>
      <div className="space-y-6">
        {comments.map((comment) => (
          <CommentCard
            key={comment.id}
            comment={comment}
            onEdit={onEditComment}
            onReply={onReply}
            onLoadReplies={onLoadReplies}
          />
        ))}
      </div>

      {loadingMore && (
        <div className="mt-8 space-y-6">
          {Array.from({ length: 3 }).map((_, index) => (
            <CommentSkeleton key={`loading-${index}`} />
          ))}
        </div>
      )}

      {meta && !meta.hasNext && comments.length > 0 && (
        <div className="text-center py-8">
          <p className="text-gray-500">You've reached the end of the comments!</p>
        </div>
      )}
    </>
  );
};

export default CommentsList;