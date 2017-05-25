import * as React from 'react';
// import Dialog from 'material-ui/Dialog';
// import FlatButton from 'material-ui/FlatButton';
import RaisedButton from 'material-ui/RaisedButton';
import LoginDialog from './LoginDialog';

interface State {
    open: boolean;
}
/**
 * A modal dialog can only be closed by selecting one of the actions.
 */
export default class DialogExampleModal extends React.Component<{}, State> {

    constructor() {
        super();
        this.state = { open: false };
    }

    handleOpen = () => {
        this.setState({ open: true });
    }

    handleClose = () => {
        this.setState({ open: false });
    }

    render() {
        /*const actions = [(
            <FlatButton
                label="Cancel"
                primary={true}
                onTouchTap={this.handleClose}
            />), (
            <FlatButton
                label="Submit"
                primary={true}
                disabled={true}
                onTouchTap={this.handleClose}
            />),
        ];*/

        return (
            <div>
                <RaisedButton label="Modal Dialog" onTouchTap={this.handleOpen} />
                <LoginDialog open={this.state.open} />
            </div>
        );
    }
}