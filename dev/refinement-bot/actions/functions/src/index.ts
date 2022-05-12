import { App, ExpressReceiver } from "@slack/bolt";
import * as functions from 'firebase-functions';
import axios from "axios";

// @ts-nocheck
console.log('panda')
// const config = functions.config();
require("dotenv").config();
// Initializes your app with your bot token and signing secret


const expressReceiver = new ExpressReceiver({
    signingSecret: process.env.SLACK_SIGNING_SECRET as string,
    endpoints: '/events',
    processBeforeResponse: true,
});

const app = new App({
    receiver: expressReceiver,
    token: process.env.SLACK_BOT_TOKEN  as string,
    processBeforeResponse: true,

});

app.command("/jason-testing", async ({  ack, say }) => {
    try {
    console.log('panda')
      await ack();
      say("Yaaay! that command works!");
    } catch (error) {
        console.log("err")
      console.error(error);
    }
});

app.action('button', async ({ ack,  say, action }) => {
    await ack();
    await say("Yaaay! that command works!");
    console.log(JSON.stringify(action), '13')
    // Update the message to reflect the action
  });

app.action('static_select-action-1', async ({ ack,  say, action: actionBase, respond, body: bodyBase }) => {
    await ack();
    // const blocks = [
    //     {
    //     "type": "image",
    //     "title": {
    //         "type": "plain_text",
    //         "text": "My message",
    //         "emoji": true
    //     },
    //     "image_url": "https://assets3.thrillist.com/v1/image/1682388/size/tl-horizontal_main.jpg",
    //     "alt_text": "my_message"
    //     }
    // ];
    const body: any = bodyBase;
    const action: any = actionBase;
    let message = body.message.blocks[2].text.text
    message = message.replace(/\n.+/, '')
    const timestamp = `\`<!date^${parseInt(action.action_ts, 10)}^{date_num} {time_secs}|Posted ${new Date().toLocaleString()}>\``
    message += `\n\`${body.user.username}\`: ${timestamp}`
    body.message.blocks[2].text.text = message
    body.message.blocks[2].accessory.initial_option = action.selected_option;
    await say("Yaaay! that command works!");

    // May redo this to traditional oauth flow
    await axios.post('https://api.github.com/repos/sourcegraph/jgornall-playground/issues/2/labels', {"labels":["bug","enhancement"]}, {
      headers: {
        //'Authorization': `token ${process.env[`SLACK_USER_${body.user.username}`]}`,
        'Authorization': `token ${process.env.REFINEMENT_BOT}`,
        'Accept': 'application/vnd.github.symmetra-preview+json'
      }
    });
    await respond({ blocks: body.message.blocks });
    // Update the message to reflect the action
});

app.action('static_select-action-2', async ({ ack, say, action, body: bodyBase, client, respond }) => {
    await ack();
   // say("Yaaay! that command works! 2");
   const body: any = bodyBase;
   console.log(body.view, '123')
    //console.log(JSON.stringify(action, '13'))
    await respond({ "blocks": [ { "type": "section", "text": { "type": "mrkdwn", "text": ":ghost:" } } ] });
    // Update the message to reflect the action
});


// Global error handler
app.error(console.log as any);

// Handle `/echo` command invocations
app.command('/echo-from-firebase', async ({ command, ack, say }) => {
    // Acknowledge command request
    await ack();

    // Requires:
    // Add chat:write scope + invite the bot user to the channel you run this command
    // Add chat:write.public + run this command in a public channel
    await say(`You said "${command.text}"`);
});

// app.error(() => {
// 	// Check the details of the error to handle cases where you should retry sending a message or stop the app
// 	console.log('wtf')
// });



// https://{your domain}.cloudfunctions.net/slack/events
exports.slack = functions.https.onRequest(expressReceiver.app);

// (async () => {
//   const port = 3000
//   // Start your app
//   await app.start(process.env.PORT || port);
//   console.log(`⚡️ Slack Bolt app is running on port ${port}!`);
// })();
