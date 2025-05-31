import { useState, useEffect, useCallback } from 'react';
import { PostsAPI } from '@/services/api';
import { Post, PaginationMeta } from '@/types';
import { useInfiniteScroll } from '@/hooks/useInfiniteScroll';

export const useDashboard = () => {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [meta, setMeta] = useState<PaginationMeta | null>(null);
  const [error, setError] = useState<string | null>(null);

  const loadPosts = useCallback(async (page: number = 1, append: boolean = false) => {
    try {
      if (page === 1) {
        setLoading(true);
        setError(null);
      } else {
        setLoadingMore(true);
      }

      const response = await PostsAPI.getPosts(page, 10);
      
      if (append) {
        setPosts(prev => [...prev, ...response.posts]);
      } else {
        setPosts(response.posts);
      }
      setMeta(response?.meta);
    } catch (err: any) {
      setError(err.message || 'Failed to load posts');
    } finally {
      setLoading(false);
      setLoadingMore(false);
    }
  }, []);

  const loadMorePosts = useCallback(() => {
    if (meta && meta.hasNext && !loadingMore) {
      loadPosts(meta.page + 1, true);
    }
  }, [meta, loadingMore, loadPosts]);

  useInfiniteScroll({
    hasNextPage: meta?.hasNext || false,
    isLoading: loadingMore,
    loadMore: loadMorePosts,
    threshold: 200
  });

  useEffect(() => {
    loadPosts(1, false);
  }, [loadPosts]);

  const handlePostCreated = (newPost: Post) => {
    setPosts(prev => [newPost, ...prev]);
    if (meta) {
      setMeta(prev => prev ? { ...prev, total: prev.total + 1 } : null);
    }
  };

  const retryLoadPosts = () => loadPosts(1, false);

  return {
    posts,
    loading,
    loadingMore,
    meta,
    error,
    handlePostCreated,
    retryLoadPosts
  };
};
