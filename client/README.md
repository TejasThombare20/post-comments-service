# Post Comments Frontend Application

A modern React-Vite frontend application for displaying posts and comments with infinite scrolling and interactive features.

## Tech Stack

- **React 18** with TypeScript
- **Vite** for fast development and building
- **Tailwind CSS v4** for styling
- **shadcn/ui** for UI components
- **React Router** for navigation
- **Axios** for API calls

## Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── ui/             # shadcn/ui components
│   └── Navigation.tsx  # Top navigation component
├── pages/              # Page components
│   └── Home.tsx        # Landing page
├── handlers/           # API handlers
│   └── api-handlers.ts # Axios configuration and API methods
├── types/              # TypeScript type definitions
│   └── index.ts        # Post and Comment interfaces
├── utils/              # Utility functions
│   └── cn.ts           # Class name utility
├── hooks/              # Custom React hooks (to be added)
├── App.tsx             # Main application component
├── main.tsx            # Application entry point
└── index.css           # Global styles with Tailwind
```

## Features Implemented

### 1. Home Page (Landing Page)
- Large headline with gradient background
- "Get Started" button that redirects to dashboard
- Feature showcase cards
- Responsive design with modern UI

### 2. Navigation
- Horizontal navigation panel
- Active route highlighting
- Responsive design

### 3. API Handler
- Centralized Axios configuration
- Error handling
- TypeScript support
- Environment-based URL configuration

## Environment Variables

Create a `.env` file in the client directory:

```env
VITE_LOCAL_SERVER=http://localhost:3000
VITE_PROD_SERVER=https://your-production-server.com
```

## Installation & Setup

1. Install dependencies:
```bash
npm install
```

2. Start development server:
```bash
npm run dev
```

3. Build for production:
```bash
npm run build
```

## Next Steps

The following features will be implemented in the next iteration:

1. **Dashboard Page** - Display posts with infinite scrolling
2. **Post Details Page** - Show individual post with comments
3. **Comment System** - Interactive comments with replies
4. **Infinite Scrolling** - Optimized loading experience
5. **Loading States** - Proper loading indicators

## API Integration

The application is configured to work with a backend server. The API handler supports:

- GET requests for fetching posts and comments
- POST requests for creating comments
- PUT/PATCH requests for editing comments
- DELETE requests for removing comments

## Styling

The application uses Tailwind CSS v4 with a custom design system including:

- CSS custom properties for theming
- Responsive design patterns
- Modern UI components from shadcn/ui
- Smooth animations and transitions
