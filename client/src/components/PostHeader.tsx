import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import RichTextEditor from '@/components/RichTextEditor';
import RichTextDisplay from '@/components/RichTextDisplay';
import { Post } from '@/types';
import { formatDistanceToNow } from 'date-fns';

interface PostHeaderProps {
  post: Post;
  isPostOwner: boolean;
  isEditingPost: boolean;
  editPostTitle: string;
  editPostContent: string;
  onEditPostClick: () => void;
  onEditPostTitle: (title: string) => void;
  onEditPostContent: (content: string) => void;
  onSavePost: () => void;
  onCancelEdit: () => void;
}

const PostHeader: React.FC<PostHeaderProps> = ({
  post,
  isPostOwner,
  isEditingPost,
  editPostTitle,
  editPostContent,
  onEditPostClick,
  onEditPostTitle,
  onEditPostContent,
  onSavePost,
  onCancelEdit,
}) => {
  const formatDate = (dateString: string) => {
    try {
      return formatDistanceToNow(new Date(dateString), { addSuffix: true });
    } catch {
      return 'Unknown time';
    }
  };

  return (
    <Card className="mb-8 border border-gray-200">
      <CardHeader>
        {isEditingPost ? (
          <div className="space-y-4">
            <Input
              type="text"
              value={editPostTitle}
              onChange={(e) => onEditPostTitle(e.target.value)}
              className="w-full text-2xl font-bold text-gray-900 border-none outline-none bg-transparent"
              placeholder="Post title..."
            />
            <div className="flex items-center gap-2 text-sm text-gray-500">
              <span>By {post.author.display_name || post.author.username}</span>
              <span>•</span>
              <span>{formatDate(post.created_at)}</span>
              {post.commentsCount !== undefined && (
                <Badge variant="secondary" className="ml-2">
                  {post.commentsCount} comments
                </Badge>
              )}
            </div>
          </div>
        ) : (
          <>
            <div className="flex justify-between items-start">
              <CardTitle className="text-2xl font-bold text-gray-900">
                {post.title}
              </CardTitle>
              {isPostOwner && (
                <Button
                  onClick={onEditPostClick}
                  variant="ghost"
                  size="sm"
                  className="text-gray-600 hover:text-blue-600"
                >
                  <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                  Edit
                </Button>
              )}
            </div>
            <div className="flex items-center gap-2 text-sm text-gray-500">
              <span>By {post.author.display_name || post.author.username}</span>
              <span>•</span>
              <span>{formatDate(post.created_at)}</span>
              {post.commentsCount !== undefined && (
                <Badge variant="secondary" className="ml-2">
                  {post.commentsCount} comments
                </Badge>
              )}
            </div>
          </>
        )}
      </CardHeader>
      <CardContent>
        {isEditingPost ? (
          <div className="space-y-4">
            <RichTextEditor
              value={editPostContent}
              onChange={onEditPostContent}
              placeholder="Write your post content..."
              className="min-h-[200px]"
            />
            <div className="flex gap-2">
              <Button onClick={onSavePost} disabled={!editPostTitle.trim() || !editPostContent.trim()}>
                Save Changes
              </Button>
              <Button onClick={onCancelEdit} variant="outline">
                Cancel
              </Button>
            </div>
          </div>
        ) : (
          <RichTextDisplay content={post.content} className="text-gray-700 leading-relaxed" />
        )}
      </CardContent>
    </Card>
  );
};

export default PostHeader;