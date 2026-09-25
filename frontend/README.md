# Conspiracy Board Frontend

A modern React-based frontend for the Conspiracy Board graph visualization system. Built with Vite for fast development and production builds.

## Features

- **Interactive Graph Visualization**: Real-time visualization of the directed graph using vis-network
- **Add People**: Create new nodes in the graph with names and descriptions
- **Add Connection Types**: Define custom relationship types (knows, works-with, etc.)
- **Add Connections**: Create directed edges between people with optional weights
- **Auto-Refresh**: Graph updates automatically when new data is added
- **Responsive Design**: Works on desktop with clean, modern UI

## Tech Stack

- **React 18** - UI library
- **Vite** - Build tool and dev server
- **vis-network** - Graph visualization
- **Axios** - HTTP client for API calls
- **CSS3** - Modern styling with gradients and animations

## Getting Started

### Prerequisites
- Node.js 16+
- npm or yarn

### Installation

```bash
cd frontend
npm install
```

### Development

Start the development server:

```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`

**Note**: Make sure the backend is running on `http://localhost:8080`

### Build for Production

```bash
npm run build
```

Output will be in the `dist/` directory.

### Preview Production Build

```bash
npm run preview
```

## Project Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── AddPerson.jsx           # Form to add people
│   │   ├── AddConnection.jsx       # Form to add connections
│   │   ├── AddConnectionType.jsx   # Form to add connection types
│   │   ├── GraphVisualization.jsx  # Graph visualization component
│   │   └── *.css                   # Component styles
│   ├── api.js                      # API client with axios
│   ├── App.jsx                     # Main app component
│   ├── App.css                     # App styles
│   └── main.jsx                    # React entry point
├── index.html                      # HTML template
├── vite.config.js                  # Vite configuration
├── package.json                    # Dependencies
└── README.md                       # This file
```

## API Integration

The frontend connects to the backend API at `http://localhost:8080`. All API calls are made through the `api.js` module using axios.

### API Endpoints Used

- `POST /api/connection-types` - Create connection type
- `GET /api/connection-types` - List all types
- `POST /api/people` - Create person
- `GET /api/people` - List all people
- `POST /api/connections` - Create connection
- `GET /api/connections` - List all connections
- `GET /api/graph` - Get full graph

## Component Details

### AddConnectionType
Form to create new connection types. Required: type name. Optional: description.

### AddPerson
Form to create new people nodes. Required: name. Optional: description.

### AddConnection
Form to create directed edges between people. Requires:
- From Person (select dropdown)
- To Person (select dropdown)
- Connection Type (select dropdown)
- Optional: Description
- Optional: Weight (0-1)

### GraphVisualization
Interactive graph display using vis-network. Features:
- Pan and zoom with mouse
- Drag nodes to rearrange
- Click nodes/edges for more info
- Physics simulation for automatic layout
- Refresh button to reload graph

## Styling

The app uses a modern gradient color scheme:
- Primary: Purple gradient (#667eea to #764ba2)
- Accents: Green nodes, red edges
- Responsive layout with sidebar and main content

## Development Tips

- Hot Module Replacement (HMR) is enabled for fast development
- Browser dev tools work seamlessly with Vite
- API calls are logged to the console
- Errors are displayed in the UI

## Browser Support

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)

## Troubleshooting

**Frontend won't connect to backend?**
- Ensure backend is running on port 8080
- Check browser console for CORS errors
- Verify API URLs in `src/api.js`

**Graph not displaying?**
- Check backend has data
- Verify vis-network dependency is installed
- Check browser console for errors

**Styles look broken?**
- Clear browser cache (Ctrl+Shift+Delete)
- Restart dev server (`npm run dev`)

## Future Enhancements

- Delete/edit existing people and connections
- Search and filter
- Export graph as image
- Dark mode
- Advanced filtering options
- Real-time WebSocket updates
