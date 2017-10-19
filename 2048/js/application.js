animationDelay = 0;
maxScore=-1;
cacheEnabled=false;
recordUrl="http://127.0.0.1:9001/record";
computeUrl="http://127.0.0.1:9001/compute";

// Wait till the browser is ready to render the game (avoids glitches)
window.requestAnimationFrame(function () {
  new GameManager(4, KeyboardInputManager, HTMLActuator, LocalStorageManager);
});
