import { Post } from '@/types';
import PostCard from '@/components/PostCard';
import PostSkeleton from './PostSkeleton';

interface PostsListProps {
  posts: Post[];
  loadingMore: boolean;
}

const PostsList = ({ posts, loadingMore }: PostsListProps) => {
  return (
    <>
      <div className="space-y-6">
        {posts.map((post) => (
          <PostCard key={post.id} post={post} />
        ))}
      </div>

      {loadingMore && (
        <div className="mt-8 space-y-6">
          {Array.from({ length: 3 }).map((_, index) => (
            <PostSkeleton key={`loading-${index}`} />
          ))}
        </div>
      )}
    </>
  );
};

export default PostsList;