int mazeWidth;
int mazeHeight;

final int cellSize = 20;

int [][] prev;
int [][] dirSeq = {{0, 1, 2, 3},
                   {1, 0, 2, 3},
                   {0, 1, 3, 2},
                   {1, 0, 3, 2},
                   {2, 3, 0, 1},
                   {2, 3, 1, 0},
                   {3, 2, 0, 1},
                   {3, 2, 1, 0}};
                   
int currX, currY;

void initGraph() {
  prev = new int[mazeWidth][];
  for(int i = 0; i < mazeWidth; ++i) {
    prev[i] = new int[mazeHeight];
    for (int j = 0; j < mazeHeight; ++j) {
      prev[i][j] = -1;
    }
  }
}

void drawGrid(int width, int height) {
  strokeWeight(10);
  for (int i = 0; i < width; i += cellSize) {
    line(i, 0, i, height);
  }
  line(width, 0, width, height);
  for (int j = 0; j < height; j += cellSize) {
    line(0, j, width, j);
  }
  line(0, height, width, height);
}

void setup() {
  size(1000, 500);
  mazeWidth = width / cellSize;
  mazeHeight = height / cellSize;
  drawGrid(width, height);
  initGraph();
  frameRate(10);
  strokeWeight(cellSize / 8);
}

boolean expand(int x, int y) {
  fill(128, 128, 128);
  rect(x * cellSize + cellSize / 4,
       y * cellSize + cellSize / 4,
       cellSize / 2,
       cellSize / 2);
  int[] seq = dirSeq[(int)random(8)];
  int nextDir = -1;
  for (int d : seq) {
    int nextX = x;
    int nextY = y;
    switch (d) {
      case 0: nextX = x + 1; break;
      case 1: nextY = y + 1; break;
      case 2: nextX = x - 1; break;
      case 3: nextY = y - 1; break;
    }
    if (nextX >= 0 && nextX < mazeWidth &&
        nextY >= 0 && nextY < mazeHeight &&
        prev[nextX][nextY] == -1) {
      nextDir = d;
      break;
    }
  }
  
  if (nextDir == -1) {
    // All neighbors are occupied, retact.
    nextDir = prev[x][y];
    if (nextDir == -1) {
      // at the root already
      return false;
    }
    switch (nextDir) {
      case 0: currX = x + 1; break;
      case 1: currY = y + 1; break;
      case 2: currX = x - 1; break;
      case 3: currY = y - 1; break;
    }
    return true;
  }
  int prevDir = -1;
  switch (nextDir) {
    case 0: currX = x + 1; prevDir = 2; break;
    case 1: currY = y + 1; prevDir = 3; break;
    case 2: currX = x - 1; prevDir = 0; break;
    case 3: currY = y - 1; prevDir = 1; break;
  }
  stroke(255);
  strokeWeight(cellSize / 4);
  line(x * cellSize + cellSize / 2, y * cellSize + cellSize / 2,
       currX * cellSize + cellSize / 2, currY * cellSize + cellSize / 2);
  prev[currX][currY] = prevDir;
  fill(255, 128, 128);
  rect(x * cellSize + cellSize / 4,
       y * cellSize + cellSize / 4,
       cellSize / 2,
       cellSize / 2);
  return true;
}

void draw() {
  expand(currX, currY); 
}