import {useEffect, useRef, useState} from 'react';
import {Network} from 'vis-network';
import {DataSet} from 'vis-data';
import {apiClient} from '../api';
import './GraphVisualization.css';

export function GraphVisualization() {
    const networkRef = useRef(null);
    const networkInstanceRef = useRef(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [selectedEdgeId, setSelectedEdgeId] = useState(null);
    const [selectedNodeId, setSelectedNodeId] = useState(null);
    const [connectionTypes, setConnectionTypes] = useState([]);

    useEffect(() => {
        loadGraph();
        return () => {
            if (networkInstanceRef.current) {
                networkInstanceRef.current.destroy();
                networkInstanceRef.current = null;
            }
        };
    }, []);

    const loadGraph = async () => {
        setLoading(true);
        setError('');
        setSelectedEdgeId(null);
        setSelectedNodeId(null);
        try {
            const response = await apiClient.getGraph();
            const graph = response.data;
            visualizeGraph(graph);
        } catch (err) {
            console.error('Graph loading error:', err);
            setError('Failed to load graph: ' + (err.message || 'Unknown error'));
        } finally {
            setLoading(false);
        }
    };

    const handleRefresh = async () => {
        await loadGraph();
    };

    const visualizeGraph = (graph) => {
        const nodes = new DataSet();
        const edges = new DataSet();

        // Store connection types for legend
        if (graph.types && Array.isArray(graph.types)) {
            setConnectionTypes(graph.types);
        }

        // Create a map of type IDs to colors for quick lookup
        const typeColorMap = {};
        if (graph.types && Array.isArray(graph.types)) {
            graph.types.forEach((type) => {
                typeColorMap[type.id] = type.color || '#666';
            });
        }

        // Add nodes
        if (graph.nodes && Array.isArray(graph.nodes)) {
            graph.nodes.forEach((node) => {
                nodes.add({
                    id: node.person.id,
                    label: node.person.name,
                    title: node.person.description || node.person.name,
                    color: '#4CAF50',
                    font: {color: '#fff'},
                });
            });
        }

        // Collect unique connections (each is listed under both endpoints)
        const connectionsById = new Map();
        if (graph.nodes && Array.isArray(graph.nodes)) {
            graph.nodes.forEach((node) => {
                (node.connections || []).forEach((conn) => connectionsById.set(conn.id, conn));
            });
        }
        const allConnections = [...connectionsById.values()];

        // The backend guarantees at most one connection per direction (A->B),
        // so a person pair has at most two connections: A->B and B->A.
        const byDirection = new Map(
            allConnections.map((conn) => [`${conn.from_person_id}-${conn.to_person_id}`, conn])
        );

        const makeEdge = (conn, {arrow, curved}) => {
            const color = typeColorMap[conn.type_id] || '#666';
            return {
                id: conn.id,
                from: conn.from_person_id,
                to: conn.to_person_id,
                title: conn.description || `Connection ${conn.id}`,
                width: 2,
                color,
                arrows: {to: arrow ? {enabled: true, scaleFactor: 0.5} : {enabled: false}},
                smooth: curved ? {type: 'curvedCW', roundness: 0.05} : false,
                highlight: {color, width: 4},
            };
        };

        const handled = new Set();
        allConnections.forEach((conn) => {
            if (handled.has(conn.id)) return;
            handled.add(conn.id);

            const isSelfLoop = conn.from_person_id === conn.to_person_id;
            const reverse = isSelfLoop
                ? undefined
                : byDirection.get(`${conn.to_person_id}-${conn.from_person_id}`);

            if (!reverse) {
                edges.add(makeEdge(conn, {arrow: true, curved: false}));
                return;
            }

            handled.add(reverse.id);
            if (conn.type_id === reverse.type_id) {
                // Same type in both directions: show as a single undirected edge
                edges.add(makeEdge(conn, {arrow: false, curved: false}));
            } else {
                // Different types: show both directions as curved arrows
                edges.add(makeEdge(conn, {arrow: true, curved: true}));
                edges.add(makeEdge(reverse, {arrow: true, curved: true}));
            }
        });

        const data = {nodes, edges};
        const options = {
            physics: {
                enabled: false,
            },
            interaction: {
                navigationButtons: true,
                keyboard: true,
                dragNodes: true,
            },
            nodes: {
                shape: 'box',
                margin: 10,
                widthConstraint: {maximum: 200},
            },
            edges: {
                smooth: {
                    type: 'continuous',
                },
            },
        };

        if (networkRef.current) {
            // Destroy previous network instance
            if (networkInstanceRef.current) {
                networkInstanceRef.current.destroy();
            }
            
            networkInstanceRef.current = new Network(networkRef.current, data, options);

            // Handle node selection
            networkInstanceRef.current.on('selectNode', (event) => {
                if (event.nodes && event.nodes.length > 0) {
                    setSelectedNodeId(event.nodes[0]);
                    setSelectedEdgeId(null);
                }
            });

            // Handle edge selection
            networkInstanceRef.current.on('selectEdge', (event) => {
                if (event.edges && event.edges.length > 0) {
                    setSelectedEdgeId(event.edges[0]);
                    setSelectedNodeId(null);
                }
            });

            // Handle deselection
            networkInstanceRef.current.on('deselectNode', () => {
                setSelectedNodeId(null);
            });

            networkInstanceRef.current.on('deselectEdge', () => {
                setSelectedEdgeId(null);
            });

            networkInstanceRef.current.on('click', (event) => {
                if ((!event.edges || event.edges.length === 0) && (!event.nodes || event.nodes.length === 0)) {
                    setSelectedEdgeId(null);
                    setSelectedNodeId(null);
                }
            });
        }
    };

    const handleDeleteConnection = async () => {
        if (!selectedEdgeId) return;

        if (!window.confirm('Delete this connection?')) {
            return;
        }

        try {
            await apiClient.deleteConnection(selectedEdgeId);
            setSelectedEdgeId(null);
            loadGraph();
        } catch (err) {
            setError('Failed to delete connection');
            console.error(err);
        }
    };

    const handleDeletePerson = async () => {
        if (!selectedNodeId) return;

        if (!window.confirm('Delete this person?')) {
            return;
        }

        try {
            await apiClient.deletePerson(selectedNodeId);
            setSelectedNodeId(null);
            loadGraph();
        } catch (err) {
            setError('Failed to delete person');
            console.error(err);
        }
    };

    return (
        <div className="graph-container">
            <div className="graph-header">
                <h2>Graph Visualization</h2>
                {selectedEdgeId && (
                    <button className="delete-btn" onClick={handleDeleteConnection}>
                        Delete Connection
                    </button>
                )}
                {selectedNodeId && (
                    <button className="delete-btn" onClick={handleDeletePerson}>
                        Delete Person
                    </button>
                )}
            </div>

            {loading && <p>Loading graph...</p>}
            {error && <p className="error">{error}</p>}
            <button onClick={handleRefresh} className="refresh-btn">
                Refresh
            </button>

            {connectionTypes.length > 0 && (
                <div className="legend">
                    <h3>Connection Types</h3>
                    <div className="legend-items">
                        {connectionTypes.map((type) => (
                            <div key={type.id} className="legend-item">
                                <div 
                                    className="legend-color" 
                                    style={{ backgroundColor: type.color || '#666' }}
                                ></div>
                                <span>{type.name}</span>
                            </div>
                        ))}
                    </div>
                </div>
            )}

            <div ref={networkRef} className="network"/>
        </div>
    );
}
