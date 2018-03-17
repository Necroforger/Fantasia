
shortCommands = {
  "chen": "honk honk",
  "ur mom gay": "no u",
};


function runCommand(name, reply) {
  for (i in shortCommands) {
    if (i === name) {
      reply(shortCommands[i]);
    }
  }
} 

// Method called on discord MessageCreate events
function onMessage(sys, msg) {

  // Shortcut function for replies
  function reply(data) {
    sys.Dream.DG.ChannelMessageSend(msg.ChannelID, data);
  }

  // Prevent the bot from responding to itself
  if (sys.Dream.DG.State.User.ID == msg.Author.ID) {
    return;
  }

  runCommand(msg.Content, reply);
}

// Method called when script is loaded
function onLoad(sys) {
  addCommand("say", "tells the bot to say something", function(ctx) {
    ctx.Reply(ctx.Args.After());
  });
}