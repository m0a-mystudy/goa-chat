import * as injectTapEventPlugin from 'react-tap-event-plugin';
injectTapEventPlugin();

import * as React from 'react';
import * as ReactDOM from 'react-dom';
import {MuiThemeProvider} from 'material-ui/styles';

import OldApp from './App';

// import RaisedButton from 'material-ui/RaisedButton';
// import {RaisedButton} from 'material-ui';
// const MyAwesomeReactComponent = () => (
//   <RaisedButton label="Default" />
// );

// import './index.css';



const App = () => (
    <MuiThemeProvider>
        <OldApp />
    </MuiThemeProvider>
);


ReactDOM.render(
    <App />,
    document.getElementById('root') as HTMLElement
);
