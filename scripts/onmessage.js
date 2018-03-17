
shortCommands = {
  "chen": "honk honk",
  "it'?s magic": "https://img.fireden.net/a/image/1464/23/1464233751072.jpg",
  "(your|ur) m[u|o]m gay": "no u",
};

function runCommand(name, reply) {
  for (i in shortCommands)
    if (RegExp("(^|\\s)" + i + "($|\\s)", "i").exec(name.toLowerCase()))
      reply(shortCommands[i]);
}

// MessageCreate event handler
function onMessage(sys, msg) {

  // Shortcut function for replies
  function reply(data) {
    sys.Dream.DG.ChannelMessageSend(msg.ChannelID, data);
  }

  // Prevent the bot from responding to itself
  if (sys.Dream.DG.State.User.ID == msg.Author.ID) return;

  runCommand(msg.Content, reply);
}

// Script load handler
function onLoad(sys) {
  addCommand("say", "tells the bot to say something", function (ctx) {
    ctx.Reply(ctx.Args.After());
  });
  addCommand("sayd", "tells the bot to say something, then deletes your message", function(ctx) {
    ctx.Reply(ctx.Args.After());
    ctx.Ses.DG.ChannelMessageDelete(ctx.Msg.ChannelID, ctx.Msg.ID);
  });
}