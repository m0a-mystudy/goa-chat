import * as React from 'react';
import './App.css';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import Room from './Room';
import Chat from './Chat';

class App extends React.Component<{}, {}> {
    render() {
        return (
            <Router>
                <div className="App">
                    <Route exact={true} path={`/`} component={Room} />
                    <Route path={`/room/:roomID`} component={Chat} />
                </div>
            </Router>
        );
    }
}

export default App;
