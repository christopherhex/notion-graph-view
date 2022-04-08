import { h } from 'preact';
import Graph from '../../components/graph';
import style from './style.css';

const Home = (props) => {

	const notionKey = props["notionKey"];

	return (
		<div class={style.home} >
			<Graph notionKey={notionKey} />
		</div>
	)
};

export default Home;
