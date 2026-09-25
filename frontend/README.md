# Conspiracy Board Frontend

A modern React-based frontend for the Conspiracy Board graph visualization system. Built with Vite for fast development and production builds.

## Features

- **Interactive Graph Visualization**: Real-time visualization of the directed graph using vis-network
- **Add People**: Create new nodes in the graph with names and descriptions
- **Manage Connection Types**: Create, edit and delete relationship types (knows, works-with, etc.) with a color each
- **Add Connections**: Create directed edges between people
- **Delete**: Select a person or connection in the graph to delete it (deleting a person also deletes their connections)
- **Auto-Refresh**: Graph and forms update automatically when data is added, edited or deleted
- **Responsive Design**: Works on desktop with clean, modern UI

## Tech Stack

- **React 19** - UI library
- **Vite 8** - Build tool and dev server
- **vis-network** - Graph visualization
- **Axios** - HTTP client for API calls
- **CSS3** - Modern styling with gradients and animations

## Getting Started

### Prerequisites
- Node.js 20.19+ or 22.12+ (required by Vite 8)
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

**Note**: Make sure the backend is running on `http://localhost:8000`

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
│   │   ├── ManageConnectionTypes.jsx # Create/edit/delete connection types
│   │   ├── Modal.jsx               # Generic modal dialog
│   │   ├── GraphVisualization.jsx  # Graph visualization component
│   │   └── *.css                   # Component styles (FormCard.css is shared by the forms)
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

All API calls are made through the `api.js` module using axios. The base URL is chosen as follows:

1. `VITE_API_URL` if set at build time
2. `http://localhost:8000/api` when the page is served from `localhost` on a port other than 8000 (e.g. the Vite dev server)
3. otherwise `<current origin>/api` (the Docker image, where nginx proxies `/api`)

Error responses have the shape `{"error": "<message>"}`; the message is shown in the UI.

### API Endpoints Used

- `POST /api/connection-types` - Create connection type
- `GET /api/connection-types` - List all types
- `PUT /api/connection-types/:id` - Update connection type
- `DELETE /api/connection-types/:id` - Delete connection type
- `POST /api/people` - Create person
- `GET /api/people` - List all people
- `DELETE /api/people/:id` - Delete person
- `POST /api/connections` - Create connection
- `DELETE /api/connections/:id` - Delete connection
- `GET /api/graph` - Get full graph

## Component Details

### Modal
Accessible dialog (`role="dialog"`, labelled by its title). Focus moves into the dialog when it opens, stays trapped inside while open and returns to the opening button on close; the rest of the page is `inert` meanwhile. Close with Escape, the × button or a click on the backdrop.

### ManageConnectionTypes
Modal content to create, edit and delete connection types. Required: type name. Optional: description, color.

### AddPerson
Form to create new people nodes. Required: name. Optional: description.

### AddConnection
Form to create directed edges between people. Requires:
- From Person (select dropdown)
- To Person (select dropdown)
- Connection Type (select dropdown)
- Optional: Description

### GraphVisualization
Interactive graph display using vis-network. Features:
- Pan and zoom with mouse
- Drag nodes to rearrange
- Hover nodes/edges to see their description
- Click a node or edge to select it for deletion
- Edges are colored by connection type (legend above the graph); a pair connected in both directions with the same type is drawn as one undirected edge
- Refresh button to reload graph

## Styling

The app uses a modern gradient color scheme:
- Primary: Purple gradient (#667eea to #764ba2)
- Accents: Green nodes, edges colored by connection type
- Responsive layout with sidebar and main content

## Development Tips

- Hot Module Replacement (HMR) is enabled for fast development
- Browser dev tools work seamlessly with Vite
- Failed API calls are logged to the console
- Errors are displayed in the UI

## Browser Support

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)

## Troubleshooting

**Frontend won't connect to backend?**
- Ensure backend is running on port 8000
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

- Edit existing people and connections
- Search and filter
- Export graph as image
- Dark mode
- Advanced filtering options
- Real-time WebSocket updates
