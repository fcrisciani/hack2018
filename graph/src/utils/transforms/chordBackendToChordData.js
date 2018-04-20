import cloneDeep from 'lodash/cloneDeep';

const separateByConnection = (matrix, labels) => {
  const connectedMatrix = [];
  const connectedLabels = [];
  const disconnectedMatrix = [];
  const disconnectedLabels = [];
  const connectedIndeces = [];
  const disconnectedIndeces = [];
  const hasConnection = (vector, matrix, i) => {
    const connectionCount = vector.reduce((a,b) => a + b, 0);
    if (connectionCount) {
      return true;
    }
    for (let j = 0; j < matrix.length; j += 1) {
      if (matrix[j][i]) {
        return true;
      }
    }
    return false;
  }
  matrix.forEach((vector, i) => {
    if (hasConnection(vector, matrix, i)) {
      connectedIndeces.push(i);
      connectedLabels.push(labels[i])
      connectedMatrix.push(vector);
    } else {
      disconnectedIndeces.push(i);
      disconnectedLabels.push(labels[i])
      disconnectedMatrix.push(vector)
    }
  })
  disconnectedIndeces.reverse().forEach(disconnectedIndex => {
    connectedMatrix.forEach(vector => {
      vector.splice(disconnectedIndex, 1)
    })
    // connectedMatrix.splice(disconnectedIndex,1)
  })
  for (let i = 0; i < connectedMatrix.length; i += 1) {
    if (connectedMatrix[i].length !== connectedMatrix.length) {
      console.error('fail assertion, connectedMatrix items should be a square matrix');
    }
  }
  return {
    connected: {
      matrix: connectedMatrix,
      labels: connectedLabels,
    },
    disconnected: {
      labels: disconnectedLabels,
    }
  }
}

export default (graph) => {
  const matrix = [];
  const labels = [];
  const disconnected = [];


  graph.forEach((serviceDesc) => {
    matrix.push([]);
    labels.push(serviceDesc.name);
  });
  graph.forEach((serviceDesc, i) => {
    serviceDesc.connections.forEach((connection, j) => {
      matrix[i][j] = connection.total;
    });
  })
  const separated = separateByConnection(cloneDeep(matrix), cloneDeep(labels));
  return {
    matrix,
    labels,
    connected: separated.connected,
    disconnected: separated.disconnected,
  }
}
