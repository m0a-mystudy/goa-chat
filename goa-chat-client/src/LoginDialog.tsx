
import * as React from 'react';
import Dialog from 'material-ui/Dialog';


interface DialogProps {
    open: boolean;
}
export default class DialogExampleModal extends React.Component<DialogProps, {}> {

    constructor(props: DialogProps) {
        super(props);
    }
    render() {
        if (this.props.open) {
            return (
                <Dialog open={this.props.open}>
                    <iframe  style={{width: '100%'}} src="http://oauth.local.com:8080/login" />
                </Dialog>
            );
        }
        return null;
    }
}