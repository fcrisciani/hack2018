export default (graph) => {
  const matrix = [];
  const labels = []
  graph.forEach((serviceDesc) => {
    matrix.push([]);
    labels.push(serviceDesc.name);
  });
  graph.forEach((serviceDesc, i) => {
    serviceDesc.connections.forEach((connection, j) => {
      matrix[i][j] = connection.total;
    });
  })
  return {
    matrix,
    labels
  }
}
