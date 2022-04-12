import { forceCenter } from 'd3';
import { h } from 'preact';
import { Link } from 'preact-router/match';
import eventBus from '../../lib/eventBus';
import style from './style.css';

const Header = () => {

	const handleRefresh = () => {
		eventBus.dispatch('forceRefresh');
	}

	return (
		<header class={style.header}>
			<div class="header-left">
				<h1>Notion Graph</h1>
			</div>
			<div class="header-center">
			</div>
			<div class="header-right">
				<a href="#" onclick={handleRefresh}>Force Refresh</a>
			</div>

		</header >
	);
}

export default Header;
