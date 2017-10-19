function GameManager(size, InputManager, Actuator, StorageManager) {
  this.size           = size; // Size of the grid
  this.inputManager   = new InputManager;
  this.storageManager = new StorageManager;
  this.actuator       = new Actuator;

  this.startTiles     = 2;
  this.running        = false;
  this.recordEnabled  = false;
  this.isRecordGrids  = false;
  this.isEnableStorage = false;
  this.isAutoRestart = false;
  this.step           = 0;

  this.inputManager.on("move", this.move.bind(this));
  this.inputManager.on("restart", this.restart.bind(this));
  this.inputManager.on("keepPlaying", this.keepPlaying.bind(this));
  this.inputManager.on("hint", this.hint.bind(this));
  this.inputManager.on("autoRun", this.autoRun.bind(this));
  this.inputManager.on("recordGame", this.recordGame.bind(this));
  this.inputManager.on("enableStorage", this.enableStorage.bind(this));
  this.inputManager.on("autoRestart", this.autoRestart.bind(this));

  this.setup();
}

// Restart the game
GameManager.prototype.restart = function () {
  this.storageManager.clearGameState();
  this.actuator.continueGame(); // Clear the game won/lost message
  this.setup();
};

GameManager.prototype.hint = function () {
  this.sendPostRequest(computeUrl,{data:this.grid.toArray(),step:this.step},this.hintResult);
};

GameManager.prototype.autoRun = function () {
  if(this.running){
    this.running=false;
  }else{
    this.running=true;
    this.sendPostRequest(computeUrl,{data:this.grid.toArray(),step:this.step},this.run);
  }
  this.updateButton()
};

GameManager.prototype.updateButton = function () {
  if(!this.running){
    this.actuator.setRunButton('AutoRun');
    this.actuator.setHint('Hint');
  }else{
    this.actuator.setRunButton('Stop');
  }
};


GameManager.prototype.recordGame = function () {
  this.isRecordGrids = this.actuator.isRecordGrid()
  this.actuate();
};

GameManager.prototype.enableStorage = function () {
  this.isEnableStorage = this.actuator.isEnableStorage();
  this.actuate();
};

GameManager.prototype.autoRestart = function () {
  this.isAutoRestart = this.actuator.isAutoRestart();
  this.actuate();
};


GameManager.prototype.hintResult = function(resp) {
  var r=JSON.parse(resp.responseText);
  if(r&&r.code==0){
    this.actuator.showHint(r.data)
  }
};

// Keep playing after winning (allows going over 2048)
GameManager.prototype.keepPlaying = function () {
  this.keepPlaying = true;
  this.actuator.continueGame(); // Clear the game won/lost message
};

// Return true if the game is lost, or has won and the user hasn't kept playing
GameManager.prototype.isGameTerminated = function () {
  return this.over || (this.won && !this.keepPlaying);
};

// Set up the game
GameManager.prototype.setup = function (autoRestart) {
  var previousState = this.storageManager.getGameState();

  // Reload the game from a previous game if present
  if (previousState) {
    this.grid        = new Grid(previousState.grid.size,
                                previousState.grid.cells); // Reload grid
    this.score       = previousState.score;
    this.over        = previousState.over;
    this.won         = previousState.won;
    this.keepPlaying = previousState.keepPlaying;
    this.recordGrids = previousState.recordGrids;
    this.isRecordGrids = previousState.isRecordGrids;
    this.isEnableStorage = previousState.isEnableStorage;
    this.isAutoRestart = previousState.isAutoRestart;
    this.step = previousState.step;
    if(!this.isEnableStorage||autoRestart){
      this.grid        = new Grid(this.size);
      this.score       = 0;
      this.over        = false;
      this.won         = false;
      this.keepPlaying = false;
      this.recordGrids = new Array();
      this.step = 0;
      this.addStartTiles();
    }
  } else {
    this.grid        = new Grid(this.size);
    this.score       = 0;
    this.over        = false;
    this.won         = false;
    this.step        = 0;
    this.keepPlaying = false;
    this.running = false;
    this.recordGrids = new Array();
    this.isRecordGrids = true;
    this.isEnableStorage = true;
    this.isAutoRestart = false;
    // Add the initial tiles
    this.addStartTiles();
  }
  this.updateButton()

  // Update the actuator
  this.actuate();
};

// Set up the initial tiles to start the game with
GameManager.prototype.addStartTiles = function () {
  for (var i = 0; i < this.startTiles; i++) {
    this.addRandomTile();
  }
};

// Adds a tile in a random position
GameManager.prototype.addRandomTile = function () {
  if (this.grid.cellsAvailable()) {
    var value = Math.random() < 0.9 ? 2 : 4;
    var tile = new Tile(this.grid.randomAvailableCell(), value);

    this.grid.insertTile(tile);
  }
};

// Sends the updated grid to the actuator
GameManager.prototype.actuate = function () {
  if (this.storageManager.getBestScore() < this.score) {
    this.storageManager.setBestScore(this.score);
  }

  // Clear the state when the game is over (game over only, not win)
  if (this.over) {
    this.storageManager.clearGameState();
  } else {
    this.storageManager.setGameState(this.serialize());
  }

  this.actuator.actuate(this.grid, {
    score:      this.score,
    over:       this.over,
    won:        this.won,
    bestScore:  this.storageManager.getBestScore(),
    terminated: this.isGameTerminated(),
    isRecordGrids:this.isRecordGrids,
    isEnableStorage:this.isEnableStorage,
    isAutoRestart:this.isAutoRestart,
  });

};

// Represent the current game as an object
GameManager.prototype.serialize = function () {
  return {
    grid:        this.grid.serialize(),
    score:       this.score,
    over:        this.over,
    won:         this.won,
    keepPlaying: this.keepPlaying,
    recordGrids: this.recordGrids,
    isRecordGrids:this.isRecordGrids,
    isEnableStorage:this.isEnableStorage,
    isAutoRestart:this.isAutoRestart,
    step: this.step,
  };
};

// Save all tile positions and remove merger info
GameManager.prototype.prepareTiles = function () {
  this.grid.eachCell(function (x, y, tile) {
    if (tile) {
      tile.mergedFrom = null;
      tile.savePosition();
    }
  });
};

// Move a tile and its representation
GameManager.prototype.moveTile = function (tile, cell) {
  this.grid.cells[tile.x][tile.y] = null;
  this.grid.cells[cell.x][cell.y] = tile;
  tile.updatePosition(cell);
};

// Move tiles on the grid in the specified direction
GameManager.prototype.move = function (direction) {
  // 0: up, 1: right, 2: down, 3: left
  var self = this;

  if (this.isGameTerminated()) return; // Don't do anything if the game's over

  var cell, tile;

  var vector     = this.getVector(direction);
  var traversals = this.buildTraversals(vector);
  var moved      = false;

  // Save the current tile positions and remove merger information
  this.prepareTiles();

  // Traverse the grid in the right direction and move tiles
  traversals.x.forEach(function (x) {
    traversals.y.forEach(function (y) {
      cell = { x: x, y: y };
      tile = self.grid.cellContent(cell);

      if (tile) {
        var positions = self.findFarthestPosition(cell, vector);
        var next      = self.grid.cellContent(positions.next);

        // Only one merger per row traversal?
        if (next && next.value === tile.value && !next.mergedFrom) {
          var merged = new Tile(positions.next, tile.value * 2);
          merged.mergedFrom = [tile, next];

          self.grid.insertTile(merged);
          self.grid.removeTile(tile);

          // Converge the two tiles' positions
          tile.updatePosition(positions.next);

          // Update the score
          self.score += merged.value;

          // The mighty 2048 tile
          if(maxScore>0&&merged.value === maxScore){
             self.won = true;
          }
        } else {
          self.moveTile(tile, positions.farthest);
        }

        if (!self.positionsEqual(cell, tile)) {
          moved = true; // The tile moved from its original cell!
        }
      }
    });
  });

  if (moved) {
    this.step++;
    this.addRandomTile();

    if (!this.movesAvailable()) {
      this.over = true; // Game over!
      // Send grid 
      if(this.isRecordGrids){
        this.sendPostRequest(recordUrl,{grid:this.recordGrids,score:this.score})
      }
      if(this.isAutoRestart){
        this.setup(true);
        this.running=true;
        this.updateButton()
      }
    }

    this.actuate();
  }
  this.recordGrid()
};

// Get the vector representing the chosen direction
GameManager.prototype.getVector = function (direction) {
  // Vectors representing tile movement
  var map = {
    0: { x: 0,  y: -1 }, // Up
    1: { x: 1,  y: 0 },  // Right
    2: { x: 0,  y: 1 },  // Down
    3: { x: -1, y: 0 }   // Left
  };

  return map[direction];
};

// Build a list of positions to traverse in the right order
GameManager.prototype.buildTraversals = function (vector) {
  var traversals = { x: [], y: [] };

  for (var pos = 0; pos < this.size; pos++) {
    traversals.x.push(pos);
    traversals.y.push(pos);
  }

  // Always traverse from the farthest cell in the chosen direction
  if (vector.x === 1) traversals.x = traversals.x.reverse();
  if (vector.y === 1) traversals.y = traversals.y.reverse();

  return traversals;
};

GameManager.prototype.findFarthestPosition = function (cell, vector) {
  var previous;

  // Progress towards the vector direction until an obstacle is found
  do {
    previous = cell;
    cell     = { x: previous.x + vector.x, y: previous.y + vector.y };
  } while (this.grid.withinBounds(cell) &&
           this.grid.cellAvailable(cell));

  return {
    farthest: previous,
    next: cell // Used to check if a merge is required
  };
};

GameManager.prototype.movesAvailable = function () {
  return this.grid.cellsAvailable() || this.tileMatchesAvailable();
};

// Check for available matches between tiles (more expensive check)
GameManager.prototype.tileMatchesAvailable = function () {
  var self = this;

  var tile;

  for (var x = 0; x < this.size; x++) {
    for (var y = 0; y < this.size; y++) {
      tile = this.grid.cellContent({ x: x, y: y });

      if (tile) {
        for (var direction = 0; direction < 4; direction++) {
          var vector = self.getVector(direction);
          var cell   = { x: x + vector.x, y: y + vector.y };

          var other  = self.grid.cellContent(cell);

          if (other && other.value === tile.value) {
            return true; // These two tiles can be merged
          }
        }
      }
    }
  }

  return false;
};

GameManager.prototype.positionsEqual = function (first, second) {
  return first.x === second.x && first.y === second.y;
};

GameManager.prototype.recordGrid = function () {
  this.recordGrids.push(this.grid.toArray())
  if(this.recordGrids.length>5){
    this.recordGrids.shift()
  }
  //console.log(JSON.stringify(this.recordGrids))
}

GameManager.prototype.sendPostRequest = function (url,data,callback) {
  var self=this;
  var xhr = new XMLHttpRequest();
  xhr.open("POST", url, true);

  xhr.setRequestHeader("Content-Type", "application/json");

  xhr.onreadystatechange = function() {//Call a function when the state changes.
      if(xhr.readyState == 4 && xhr.status == 200) {
        self.func=callback;
        self.func(xhr);
      }
  }
  xhr.send(JSON.stringify(data));
}

// moves continuously until game is over
GameManager.prototype.run = function(resp) {

  var r=JSON.parse(resp.responseText);
  if(r&&r.code==0){
    if(!this.running){
      return
    }
    if(this.over){
      this.running=false;
      this.updateButton()
      return
    }

    this.actuator.showHint(r.data)
    this.move(r.data);

    var self = this;
    var timeout = animationDelay;

    if (this.running && !this.over && !this.won) {
      setTimeout(function(){
        self.sendPostRequest(computeUrl,{data:self.grid.toArray(),step:self.step},self.run);
      }, timeout);

    }
  }
}
