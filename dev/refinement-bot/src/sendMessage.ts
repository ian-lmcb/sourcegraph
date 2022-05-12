
import fetch from 'node-fetch';

import { GithubIssue, slackMessage } from './types'

const issue: string = (process.env.RANDOM_ISSUE || process.env.LABELED_ISSUE  || '{}')
const slackurl: string = process.env.SLACK_WEBHOOK_URL as string;
const checkPriority: boolean = process.env.CHECK_PRIORITY as string === 'true';
if (!issue) {
    console.assert(issue, 'issue exists',)
    process.exit()
}

const isPriorityPresent = (json_issue: GithubIssue) =>
     false
    // let priorityPresent = false;
    // priorityPresent = json_issue.labels?.some((label) => {
    //     // match p0, p1, p2
    //     return label.name.match(/^fep\/p\d$/)
    // }) || false;
    // return priorityPresent

const isEstimatePresent = (json_issue: GithubIssue) =>
     false
    // let estimatePresent = false;
    // estimatePresent = json_issue.labels?.some((label) => {
    //     return label.name.startsWith('estimate/')
    // }) || false
    // return estimatePresent

const eligibleToAlert = (json_issue: GithubIssue) => {
    const missingData = !isPriorityPresent(json_issue) || !isEstimatePresent(json_issue)
    return missingData;
}

const generateSlackTemplate = (json_issue: GithubIssue) => {
    // check estimate

    // built via https://api.slack.com/block-kit
    const slackMessage: slackMessage = {
        'blocks': [
            {
                'type': 'section',
                'text': {
                    'type': 'mrkdwn',
                    'text': `${json_issue.title} <${json_issue.html_url}|#${json_issue.number}>`
                }
            },
            {
                'type': 'section',
                'fields': []
            },
            {
                'type': 'section',
                'text': {
                    'type': 'mrkdwn',
                    'text': 'Please :thread: off this message to discuss'
                }
            }
        ]
    }

    const estimatePresent = isEstimatePresent(json_issue)
    if (!estimatePresent) {
        console.log(slackMessage)
        slackMessage.blocks![1].fields!.push(
            {
                'type': 'mrkdwn',
                'text': '*missing <https://github.com/sourcegraph/sourcegraph/labels?q=estimate|estimate>:*\n`estimate/0.5d` `estimate/5d`...'
            }
        )
    }

    // check priority
    const priorityPresent = isPriorityPresent(json_issue)
    if (!priorityPresent && checkPriority) {
        slackMessage.blocks![1].fields!.push(
            {
                'type': 'mrkdwn',
                'text': '*missing <https://github.com/sourcegraph/sourcegraph/labels?q=fep|priority>*: `fep/p0`, `fep/p1`..'
            }
        )
    }
    return slackMessage;
}

const sendMessage = async () => {
    const json_issue: GithubIssue = JSON.parse(issue)
    if (!eligibleToAlert(json_issue)){
        console.log('skipped alert')
        return;
    }
    const response = await fetch(slackurl, {
        body: JSON.stringify(generateSlackTemplate(json_issue)),
        headers: {
        'Content-Type': 'application/json'
        },
        method: 'POST'
    })
    console.log('response', response)
}

sendMessage();
