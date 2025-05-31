import { Button } from '@/components/ui/button';
import { User, PaginationMeta } from '@/types';

interface DashboardHeaderProps {
  user: User | null;
  postsCount: number;
  meta: PaginationMeta | null;
  onCreatePost: () => void;
}

const DashboardHeader = ({ user, postsCount, meta, onCreatePost }: DashboardHeaderProps) => {
  return (
    <div className="mb-8 flex justify-between items-center">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          Welcome back {user && ","} {user?.display_name || user?.username} {user && "!"}
        </h1>
        <p className="text-gray-600">  
          Discover and explore posts from the community
          {meta && (
            <span className="ml-2">
              • Showing {postsCount} of {meta.total} posts
            </span>
          )}
        </p>
      </div>
      <div className="flex gap-3">
        <Button 
          onClick={onCreatePost}
          className="flex items-center gap-2"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          Create Post
        </Button>
      </div>
    </div>
  );
};

export default DashboardHeader;
