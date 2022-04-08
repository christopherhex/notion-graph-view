import { h } from 'preact';
import Header from './header';
import Router from "preact-router"
import Home from '../routes/home';

// Code-splitting is automated for `routes` directory
import Graph from './graph';

const App = () => {

	return (
		<div id="app">
			<Header />
			<Router>
				<Home path="/:params?" ></Home>
			</Router>
		</div>
	);
}

export default App;
