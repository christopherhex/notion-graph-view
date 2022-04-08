import { h } from 'preact';
import Head from "preact-head";
import Header from './header';
import Router from "preact-router"
import Home from '../routes/home';

// Code-splitting is automated for `routes` directory
import Graph from './graph';

const App = () => {

	return (
		<div id="app">
			<Head>
				<script src="/assets/wasm_exec.js"></script>
			</Head>
			<Header />
			<Router>
				<Home path="/:params?" ></Home>
			</Router>
		</div>
	);
}

export default App;
