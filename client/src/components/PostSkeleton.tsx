import { Skeleton } from '@/components/ui/skeleton';

const PostSkeleton = () => (
  <div className="w-full p-6 border border-gray-200 rounded-lg">
    <div className="flex justify-between items-start mb-4">
      <Skeleton className="h-6 w-3/4" />
      <Skeleton className="h-5 w-20" />
    </div>
    <div className="flex items-center gap-2 mb-4">
      <Skeleton className="h-4 w-16" />
      <Skeleton className="h-4 w-2" />
      <Skeleton className="h-4 w-20" />
    </div>
    <Skeleton className="h-4 w-full mb-2" />
    <Skeleton className="h-4 w-full mb-2" />
    <Skeleton className="h-4 w-2/3 mb-4" />
    <div className="flex justify-between items-center">
      <Skeleton className="h-9 w-32" />
      <Skeleton className="h-4 w-24" />
    </div>
  </div>
);

export default PostSkeleton;