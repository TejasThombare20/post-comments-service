import PostSkeleton from './PostSkeleton';
import { Skeleton } from '@/components/ui/skeleton';

interface LoadingSkeletonsProps {
  count?: number;
  showHeader?: boolean;
}

const LoadingSkeletons = ({ count = 5, showHeader = false }: LoadingSkeletonsProps) => {
  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        {showHeader && (
          <div className="mb-8">
            <Skeleton className="h-8 w-48 mb-2" />
            <Skeleton className="h-4 w-96" />
          </div>
        )}
        <div className="space-y-6">
          {Array.from({ length: count }).map((_, index) => (
            <PostSkeleton key={index} />
          ))}
        </div>
      </div>
    </div>
  );
};

export default LoadingSkeletons;