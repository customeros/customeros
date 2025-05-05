const core = require('@actions/core');
const github = require('@actions/github');

async function run() {
  try {
    const { eventName, payload } = github.context;
    const pr = payload.pull_request;
    
    // Environment variables
    const slackApiToken = process.env.SLACK_API_TOKEN;
    const slackChannelId = process.env.SLACK_CHANNEL_ID;
    
    if (!slackApiToken || !slackChannelId) {
      throw new Error('Missing required environment variables: SLACK_API_TOKEN or SLACK_CHANNEL_ID');
    }
    
    // Constants for reactions
    const REACTIONS = {
      APPROVED: 'white_check_mark',
      MERGED: 'partymerge',
      CLOSED: 'x',
      COMMENT: 'eyes'
    };
    
    // Find existing PR message
    const existingMessage = await findPrMessage(pr.html_url, slackChannelId, slackApiToken);
    
    // Determine action type
    let action = determineAction(eventName, payload, pr);
    
    // Handle the action
    if (action === 'opened') {
      const newMessage = await postNewMessage(pr, slackChannelId, slackApiToken);
      if (newMessage) {
        await updateReactions(newMessage, action, REACTIONS, slackApiToken);
      }
    } else if (existingMessage && action) {
      await updateReactions(existingMessage, action, REACTIONS, slackApiToken);
    }
  } catch (error) {
    core.setFailed(`Action failed: ${error.message}`);
  }
}

// Determine what action to take based on the event type
function determineAction(eventName, payload, pr) {
  if (eventName === 'pull_request') {
    if (payload.action === 'opened' || payload.action === 'ready_for_review') {
      return 'opened';
    } else if (payload.action === 'synchronize') {
      return 'synchronize';
    } else if (payload.action === 'closed') {
      return pr.merged ? 'merged' : 'closed';
    }
  } else if (eventName === 'pull_request_review') {
    if (payload.review && payload.review.state === 'approved') {
      return 'approved';
    }
  } else if (eventName === 'issue_comment' && payload.issue.pull_request) {
    return 'commented';
  }
  return '';
}

// Find an existing PR message in Slack
async function findPrMessage(prUrl, channelId, token) {
  try {
    const historyResponse = await fetch(`https://slack.com/api/conversations.history?channel=${channelId}&limit=20`, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    
    const historyData = await historyResponse.json();
    if (historyData.ok && historyData.messages.length > 0) {
      const prMessage = historyData.messages.find(msg => 
        msg.text && msg.text.includes(prUrl)
      );
      
      if (prMessage) {
        return { ts: prMessage.ts, channel: channelId };
      }
    }
    return null;
  } catch (error) {
    console.error('Error searching for PR message:', error.message);
    return null;
  }
}

// Post a new message to Slack about the PR
async function postNewMessage(pr, channelId, token) {
  try {
    const userResponse = await fetch(`https://api.github.com/users/${pr.user.login}`, {
      headers: {
        'Accept': 'application/vnd.github.v3+json',
        'Authorization': `token ${process.env.GITHUB_TOKEN}`
      }
    });
    
    const userData = await userResponse.json();
    const displayName = userData.name || pr.user.login;
    const repoName = github.context.repo.repo;
    const message = `\`${repoName}\` <${pr.html_url}|${pr.title}>`;
    
    const response = await fetch('https://slack.com/api/chat.postMessage', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        channel: channelId,
        text: message,
        unfurl_links: false,
        unfurl_media: false,
        username: displayName,
        icon_url: pr.user.avatar_url
      }),
    });
    
    const data = await response.json();
    if (!data.ok) {
      throw new Error(`Failed to post message: ${data.error}`);
    }
    
    return { ts: data.ts, channel: data.channel || channelId };
  } catch (error) {
    console.error('Error posting new PR message:', error.message);
    return null;
  }
}

// Update reactions on a Slack message
async function updateReactions(messageInfo, action, REACTIONS, token) {
  try {
    // Get current reactions
    const reactionsResponse = await fetch(`https://slack.com/api/reactions.get?channel=${messageInfo.channel}&timestamp=${messageInfo.ts}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    
    const reactionsData = await reactionsResponse.json();
    const currentReactions = reactionsData.ok && reactionsData.message && reactionsData.message.reactions ? 
      reactionsData.message.reactions.map(r => r.name) : [];
    
    // Update reactions based on action
    switch(action) {
      case 'commented':
        if (!currentReactions.includes(REACTIONS.COMMENT)) {
          await addReaction(messageInfo, REACTIONS.COMMENT, token);
        }
        break;
        
      case 'approved':
        if (!currentReactions.includes(REACTIONS.APPROVED)) {
          await addReaction(messageInfo, REACTIONS.APPROVED, token);
        }
        break;
        
      case 'merged':
        if (!currentReactions.includes(REACTIONS.MERGED)) {
          await addReaction(messageInfo, REACTIONS.MERGED, token);
          // Remove other status reactions
          if (currentReactions.includes(REACTIONS.APPROVED)) {
            await removeReaction(messageInfo, REACTIONS.APPROVED, token);
          }
          if (currentReactions.includes(REACTIONS.CLOSED)) {
            await removeReaction(messageInfo, REACTIONS.CLOSED, token);
          }
          if (currentReactions.includes(REACTIONS.COMMENT)) {
            await removeReaction(messageInfo, REACTIONS.COMMENT, token);
          }
        }
        break;
        
      case 'closed':
        if (!currentReactions.includes(REACTIONS.CLOSED) && !currentReactions.includes(REACTIONS.MERGED)) {
          await addReaction(messageInfo, REACTIONS.CLOSED, token);
        }
        break;
    }
  } catch (error) {
    console.error('Error updating reactions:', error.message);
  }
}

// Add a reaction to a Slack message
async function addReaction(messageInfo, emoji, token) {
  try {
    const response = await fetch('https://slack.com/api/reactions.add', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        channel: messageInfo.channel,
        name: emoji,
        timestamp: messageInfo.ts,
      }),
    });
    
    const data = await response.json();
    if (!data.ok) {
      console.error(`Failed to add reaction: ${data.error}`);
    }
  } catch (error) {
    console.error(`Error adding reaction ${emoji}:`, error.message);
  }
}

// Remove a reaction from a Slack message
async function removeReaction(messageInfo, emoji, token) {
  try {
    const response = await fetch('https://slack.com/api/reactions.remove', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        channel: messageInfo.channel,
        name: emoji,
        timestamp: messageInfo.ts,
      }),
    });
    
    const data = await response.json();
    if (!data.ok) {
      console.error(`Failed to remove reaction: ${data.error}`);
    }
  } catch (error) {
    console.error(`Error removing reaction ${emoji}:`, error.message);
  }
}

module.exports = run;

// Execute the function if this file is run directly
if (require.main === module) {
  run();
}
