import * as injectTapEventPlugin from 'react-tap-event-plugin';
injectTapEventPlugin();
import App from './App';

import * as React from 'react';
import * as ReactDOM from 'react-dom';
import {MuiThemeProvider} from 'material-ui/styles';

const Index = () => (
    <MuiThemeProvider>
        <App />
    </MuiThemeProvider>
);

ReactDOM.render(
    <Index />,
    document.getElementById('root') as HTMLElement
);
