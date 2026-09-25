import {useEffect, useRef, useState} from 'react';
import {Network} from 'vis-network';
import {DataSet} from 'vis-data';
import {apiClient} from '../api';
import './GraphVisualization.css';

export function GraphVisualization({ onRefresh }) {
    const networkRef = useRef(null);
    const networkInstanceRef = useRef(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [selectedEdgeId, setSelectedEdgeId] = useState(null);
    const [selectedNodeId, setSelectedNodeId] = useState(null);

    useEffect(() => {
        loadGraph();
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
        if (onRefresh) {
            onRefresh();
        }
    };

    const visualizeGraph = (graph) => {
        const nodes = new DataSet();
        const edges = new DataSet();

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

        // Add edges (with deduplication)
        const addedEdges = new Set();
        const connectionPairs = new Map(); // Track bidirectional connections

        // First pass: collect all connections and detect bidirectional pairs
        const allConnections = [];
        if (graph.nodes && Array.isArray(graph.nodes)) {
            graph.nodes.forEach((node) => {
                if (node.connections && Array.isArray(node.connections)) {
                    node.connections.forEach((conn) => {
                        if (!addedEdges.has(conn.id)) {
                            addedEdges.add(conn.id);
                            allConnections.push(conn);
                        }
                    });
                }
            });
        }

        // Identify bidirectional connection pairs
        allConnections.forEach((conn) => {
            const reverse = allConnections.find(
                (c) => c.from_person_id === conn.to_person_id && c.to_person_id === conn.from_person_id
            );
            if (reverse) {
                const key = [conn.from_person_id, conn.to_person_id].sort().join('-');
                if (!connectionPairs.has(key)) {
                    connectionPairs.set(key, [conn, reverse]);
                }
            }
        });

        // Add edges with smart handling for bidirectional connections
        const processedPairs = new Set();
        allConnections.forEach((conn) => {
            const key = [conn.from_person_id, conn.to_person_id].sort().join('-');
            const pairData = connectionPairs.get(key);
            const isBidirectional = pairData && pairData.length > 1;

            if (isBidirectional) {
                // Skip if we already processed this pair
                if (processedPairs.has(key)) return;
                processedPairs.add(key);

                const [conn1, conn2] = pairData;
                const sameType = conn1.type_id === conn2.type_id;

                if (sameType) {
                    // Same type: show as single undirected edge
                    edges.add({
                        id: conn1.id, // Use first connection's ID
                        from: conn1.from_person_id,
                        to: conn1.to_person_id,
                        title: conn1.description || `Connection ${conn1.id}`,
                        width: 2,
                        color: typeColorMap[conn1.type_id] || '#666',
                        arrows: {
                            to: {enabled: false}, // No arrows for undirected
                        },
                        smooth: false,
                        highlight: {
                            color: typeColorMap[conn1.type_id] || '#666',
                            width: 4,
                        },
                    });
                } else {
                    // Different types: show both with curves
                    edges.add({
                        id: conn1.id,
                        from: conn1.from_person_id,
                        to: conn1.to_person_id,
                        title: conn1.description || `Connection ${conn1.id}`,
                        width: 2,
                        color: typeColorMap[conn1.type_id] || '#666',
                        arrows: {
                            to: {
                                enabled: true,
                                scaleFactor: 0.5,
                            },
                        },
                        smooth: {
                            type: 'curvedCW',
                            roundness: 0.2
                        },
                        highlight: {
                            color: typeColorMap[conn1.type_id] || '#666',
                            width: 4,
                        },
                    });
                    edges.add({
                        id: conn2.id,
                        from: conn2.from_person_id,
                        to: conn2.to_person_id,
                        title: conn2.description || `Connection ${conn2.id}`,
                        width: 2,
                        color: typeColorMap[conn2.type_id] || '#666',
                        arrows: {
                            to: {
                                enabled: true,
                                scaleFactor: 0.5,
                            },
                        },
                        smooth: {
                            type: 'curvedCW',
                            roundness: 0.2
                        },
                        highlight: {
                            color: typeColorMap[conn2.type_id] || '#666',
                            width: 4,
                        },
                    });
                }
            } else {
                // Unidirectional: show normally
                edges.add({
                    id: conn.id,
                    from: conn.from_person_id,
                    to: conn.to_person_id,
                    title: conn.description || `Connection ${conn.id}`,
                    width: 2,
                    color: typeColorMap[conn.type_id] || '#666',
                    arrows: {
                        to: {
                            enabled: true,
                            scaleFactor: 0.5,
                        },
                    },
                    smooth: false,
                    highlight: {
                        color: typeColorMap[conn.type_id] || '#666',
                        width: 4,
                    },
                });
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

            <div ref={networkRef} className="network"/>
        </div>
    );
}
