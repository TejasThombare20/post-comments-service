import { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useDashboard } from '@/components/UseDashboard';
import LoadingSkeletons from '@/components/LoadingSkeleton';
import ErrorState from '@/components/DashboardErrorState';
import DashboardHeader from '@/components/DashboardHeader';
import EmptyState from '@/components/EmptyState';
import PostsList from '@/components/PostList';
import CreatePostModal from '@/components/CreatePostModal';


const Dashboard = () => {
  const { user } = useAuth();
  const [showCreateModal, setShowCreateModal] = useState(false);
  
  const {
    posts,
    loading,
    loadingMore,
    meta,
    error,
    handlePostCreated,
    retryLoadPosts
  } = useDashboard();

  const handleCreatePost = () => setShowCreateModal(true);
  const handleCloseModal = () => setShowCreateModal(false);
  
  const handlePostCreatedWithClose = (newPost: any) => {
    handlePostCreated(newPost);
    setShowCreateModal(false);
  };

  // Show loading skeleton on initial load
  if (loading && posts.length === 0) {
    return <LoadingSkeletons showHeader />;
  }

  // Show error state if there's an error and no posts
  if (error && posts.length === 0) {
    return <ErrorState error={error} onRetry={retryLoadPosts} />;
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <DashboardHeader
          user={user}
          postsCount={posts.length}
          meta={meta}
          onCreatePost={handleCreatePost}
        />

        {posts.length === 0 && !loading && !error ? (
          <EmptyState onCreatePost={handleCreatePost} />
        ) : (
          <>
            <PostsList posts={posts} loadingMore={loadingMore} />
            {meta?.hasNext || posts.length === 0 && (
            <div className="text-center py-8">
                  <p className="text-gray-500">You've reached the end of the posts!</p>
              </div>
            )}
          </>
        )}
      </div>

      <CreatePostModal
        isOpen={showCreateModal}
        onClose={handleCloseModal}
        onPostCreated={handlePostCreatedWithClose}
      />
    </div>
  );
};

export default Dashboard;
