import { useNavigate } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Post } from '@/types';
import { formatDistanceToNow } from 'date-fns';
import RichTextDisplay from '@/components/RichTextDisplay';

interface PostCardProps {
  post: Post;
}

const PostCard = ({ post }: PostCardProps) => {
  const navigate = useNavigate();

  const handleLoadComments = () => {
    navigate(`/dashboard/${post.id}`);
  };

  const formatDate = (dateString: string) => {
    try {
      return formatDistanceToNow(new Date(dateString), { addSuffix: true });
    } catch {
      return 'Unknown time';
    }
  };

  return (
    <Card className="w-full hover:shadow-lg transition-shadow duration-200 border border-gray-200">
      <CardHeader className="pb-3">
        <div className="flex justify-between items-start">
          <CardTitle className="text-xl font-semibold text-gray-900 line-clamp-2">
            {post.title}
          </CardTitle>
          {post.commentsCount !== undefined && (
            <Badge variant="secondary" className="ml-2 flex-shrink-0">
              {post.commentsCount} comments
            </Badge>
          )}
        </div>
        <div className="flex items-center text-sm text-gray-500 mt-2">
          <span>By {post.author.display_name || post.author.username}</span>
          <span className="mx-2">•</span>
          <span>{formatDate(post.created_at)}</span>
        </div>
      </CardHeader>
      
      <CardContent className="pt-0">
        <div className="mb-4 line-clamp-3">
          <RichTextDisplay content={post.content} className="text-gray-700" />
        </div>
        
        <div className="flex justify-between items-center">
          <Button
            onClick={handleLoadComments}
            variant="outline"
            className="flex items-center gap-2"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
            </svg>
            Load Comments
          </Button>
          
          <div className="text-sm text-gray-500">
            {post.updated_at !== post.created_at && (
              <span>Updated {formatDate(post.updated_at)}</span>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default PostCard; 