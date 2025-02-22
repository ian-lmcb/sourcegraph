import { storiesOf } from '@storybook/react'
import { parseISO } from 'date-fns'
import AlertCircleIcon from 'mdi-react/AlertCircleIcon'
import CheckIcon from 'mdi-react/CheckIcon'

import { Timeline } from './Timeline'
import { WebStory } from './WebStory'

const { add } = storiesOf('web/Timeline', module)

add('Empty', () => <WebStory>{() => <Timeline stages={[]} />}</WebStory>)

add('Basic', () => (
    <WebStory>
        {() => (
            <Timeline
                now={() => parseISO('2020-08-01T16:21:00+00:00')}
                stages={[
                    {
                        icon: <CheckIcon />,
                        className: 'bg-success',
                        text: 'First event description',
                        date: '2020-06-15T11:15:00+00:00',
                    },
                    {
                        icon: <AlertCircleIcon />,
                        className: 'bg-danger',
                        text: 'Second event description',
                        date: '2020-06-15T12:20:00+00:00',
                    },
                    {
                        icon: <CheckIcon />,
                        className: 'bg-success',
                        text: 'Third event description',
                        date: '2020-06-15T13:25:00+00:00',
                    },
                    {
                        icon: <AlertCircleIcon />,
                        className: 'bg-danger',
                        text: 'Fourth event description',
                        date: '2020-06-15T14:30:00+00:00',
                    },
                    {
                        icon: <CheckIcon />,
                        className: 'bg-success',
                        text: 'Fifth event description',
                        date: '2020-06-15T15:35:00+00:00',
                    },
                ]}
            />
        )}
    </WebStory>
))

add('Details', () => (
    <WebStory>
        {() => (
            <Timeline
                now={() => parseISO('2020-08-01T16:21:00+00:00')}
                stages={[
                    {
                        icon: <CheckIcon />,
                        className: 'bg-success',
                        text: 'First event description',
                        date: '2020-06-15T11:15:00+00:00',
                    },
                    {
                        icon: <AlertCircleIcon />,
                        className: 'bg-danger',
                        text: 'Second event description',
                        date: '2020-06-15T12:20:00+00:00',
                        details: <p>HELLO THERE</p>,
                        expanded: true,
                    },
                    {
                        icon: <CheckIcon />,
                        className: 'bg-success',
                        text: 'Third event description',
                        date: '2020-06-15T13:25:00+00:00',
                        details: <p>HELLO THERE</p>,
                        expanded: false,
                    },
                ]}
            />
        )}
    </WebStory>
))
