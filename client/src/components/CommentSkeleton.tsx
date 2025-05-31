import React from 'react';
import { Skeleton } from '@/components/ui/skeleton';

const CommentSkeleton: React.FC = () => (
  <div className="w-full p-4 border border-gray-200 rounded-lg">
    <div className="flex items-center gap-2 mb-3">
      <Skeleton className="h-4 w-20" />
      <Skeleton className="h-4 w-16" />
      <Skeleton className="h-4 w-12" />
    </div>
    <Skeleton className="h-4 w-full mb-2" />
    <Skeleton className="h-4 w-full mb-2" />
    <Skeleton className="h-4 w-3/4 mb-3" />
    <div className="flex gap-2">
      <Skeleton className="h-8 w-16" />
      <Skeleton className="h-8 w-16" />
      <Skeleton className="h-8 w-24" />
    </div>
  </div>
);

export default CommentSkeleton;