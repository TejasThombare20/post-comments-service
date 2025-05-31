import React from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import RichTextEditor from '@/components/RichTextEditor';
import { PaginationMeta } from '@/types';

interface CommentFormProps {
  isCommenting: boolean;
  newComment: string;
  meta: PaginationMeta | null;
  onCommentClick: () => void;
  onNewCommentChange: (value: string) => void;
  onSubmitComment: () => void;
  onCancelComment: () => void;
}

const CommentForm: React.FC<CommentFormProps> = ({
  isCommenting,
  newComment,
  meta,
  onCommentClick,
  onNewCommentChange,
  onSubmitComment,
  onCancelComment,
}) => {
  return (
    <div className="mb-6">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold text-gray-900">
          Comments
          {meta && (
            <span className="ml-2 text-sm font-normal text-gray-500">
              ({meta.total} total)
            </span>
          )}
        </h2>
        <Button onClick={onCommentClick} variant="outline">
          <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          Add Comment
        </Button>
      </div>

      {isCommenting && (
        <Card className="mb-6 border border-gray-200">
          <CardContent className="p-4">
            <div className="space-y-3">
              <RichTextEditor
                value={newComment}
                onChange={onNewCommentChange}
                placeholder="Write your comment..."
                className="min-h-[150px]"
              />
              <div className="flex gap-2">
                <Button onClick={onSubmitComment} disabled={!newComment.trim()}>
                  Post Comment
                </Button>
                <Button onClick={onCancelComment} variant="outline">
                  Cancel
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
};

export default CommentForm;