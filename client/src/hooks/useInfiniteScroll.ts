import { useState, useEffect, useCallback } from 'react';

interface UseInfiniteScrollProps {
  hasNextPage: boolean;
  isLoading: boolean;
  loadMore: () => void;
  threshold?: number;
}

export const useInfiniteScroll = ({
  hasNextPage,
  isLoading,
  loadMore,
  threshold = 100
}: UseInfiniteScrollProps) => {
  const [isFetching, setIsFetching] = useState(false);

  const handleScroll = useCallback(() => {
    if (window.innerHeight + document.documentElement.scrollTop >= 
        document.documentElement.offsetHeight - threshold) {
      if (hasNextPage && !isLoading && !isFetching) {
        setIsFetching(true);
      }
    }
  }, [hasNextPage, isLoading, isFetching, threshold]);

  useEffect(() => {
    if (isFetching && hasNextPage && !isLoading) {
      loadMore();
      setIsFetching(false);
    }
  }, [isFetching, hasNextPage, isLoading, loadMore]);

  useEffect(() => {
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, [handleScroll]);

  return { isFetching };
}; 