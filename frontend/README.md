# Read-It-Later Frontend

A modern Vue.js frontend application for the Read-It-Later service, allowing users to save, organize, and manage articles for later reading.

## Features

- **User Authentication**: Secure login and registration
- **Article Management**: Save, organize, and manage articles
- **Search & Filter**: Find articles quickly with search and filtering
- **Responsive Design**: Works on desktop and mobile devices
- **Real-time Updates**: Live updates when articles are modified
- **Modern UI**: Clean, intuitive interface built with Tailwind CSS

## Tech Stack

- **Vue 3**: Progressive JavaScript framework
- **Vite**: Fast build tool and development server
- **Pinia**: State management for Vue
- **Vue Router**: Client-side routing
- **Tailwind CSS**: Utility-first CSS framework
- **Axios**: HTTP client for API requests

## Prerequisites

- Node.js 18+ and npm
- Backend API server running

## Installation

1. **Install dependencies**:
   ```bash
   npm install
   ```

2. **Set up environment variables**:
   ```bash
   cp .env.development .env.local
   ```
   
   Edit `.env.local` and configure:
   ```
   VITE_API_URL=http://localhost:8080
   VITE_APP_NAME=Read-It-Later
   VITE_APP_VERSION=1.0.0
   ```

## Development

1. **Start the development server**:
   ```bash
   npm run dev
   ```

2. **Open your browser** and navigate to `http://localhost:5173`

## Building for Production

1. **Build the application**:
   ```bash
   npm run build
   ```

2. **Preview the production build**:
   ```bash
   npm run preview
   ```

## Project Structure

```
frontend/
├── public/                 # Static assets
├── src/
│   ├── assets/            # CSS and other assets
│   │   ├── ArticleCard.vue
│   │   └── AddArticleModal.vue
│   ├── stores/            # Pinia stores
│   │   ├── auth.js        # Authentication state
│   │   └── articles.js    # Articles state
│   ├── utils/             # Utility functions
│   │   └── api.js         # API configuration
│   ├── views/             # Page components
│   │   ├── Home.vue
│   │   ├── Login.vue
│   │   └── Dashboard.vue
│   ├── router/            # Vue Router configuration
│   ├── App.vue            # Root component
│   └── main.js            # Application entry point
├── index.html             # HTML template
├── package.json           # Dependencies and scripts
├── vite.config.js         # Vite configuration
├── tailwind.config.js     # Tailwind CSS configuration
└── postcss.config.js      # PostCSS configuration
```

## Key Components

### ArticleCard
Displays individual articles with actions for reading, favoriting, archiving, and deleting.

### AddArticleModal
Modal component for adding new articles with URL capture and metadata extraction.

### Stores

#### Auth Store (`stores/auth.js`)
- User authentication state
- Login/logout functionality
- Token management
- User profile management

#### Articles Store (`stores/articles.js`)
- Articles collection management
- CRUD operations for articles
- Search and filtering
- Statistics tracking

## API Integration

The frontend communicates with the backend API using Axios with automatic:
- Authentication token injection
- Token refresh on expiry
- Error handling and retry logic

## Styling

The application uses Tailwind CSS for styling with:
- Responsive design patterns
- Custom color palette
- Consistent spacing and typography
- Smooth animations and transitions

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_URL` | Backend API URL | `http://localhost:8080` |
| `VITE_APP_NAME` | Application name | `Read-It-Later` |
| `VITE_APP_VERSION` | Application version | `1.0.0` |

## Browser Support

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License. 