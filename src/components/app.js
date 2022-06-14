import { h } from 'preact';
import Header from './header';
import Router from "preact-router"
import Home from '../routes/home';

// Code-splitting is automated for `routes` directory
import Graph from './graph';

const App = () => {

	return (
		<div id="app">
			<Header class="app-header" />
			<Router>
				<Home class="app-content" path="/:params?" ></Home>
			</Router>
		</div>
	);
}

export default App;
