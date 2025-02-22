import { storiesOf } from '@storybook/react'

import { WebStory } from '../../../../components/WebStory'

import { DownloadSpecModal } from './DownloadSpecModal'

const { add } = storiesOf('web/batches/batch-spec/edit/DownloadSpecModal', module).addDecorator(story => (
    <div className="p-3 container">{story()}</div>
))

add('Download Spec Modal', () => (
    <WebStory>
        {props => (
            <DownloadSpecModal
                name=""
                originalInput=""
                setIsDownloadSpecModalOpen={function (): void {
                    throw new Error('Function not implemented.')
                }}
                setDownloadSpecModalDismissed={function (): void {
                    throw new Error('Function not implemented.')
                }}
                {...props}
            />
        )}
    </WebStory>
))
