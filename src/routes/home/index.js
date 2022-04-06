import { h } from 'preact';
import Graph from '../../components/graph';
import style from './style.css';

const Home = (props) => {

	console.log(props);
	const notionKey = props["notionKey"];

	return (
		<div class={style.home} >
			<Graph notionKey={notionKey} />
		</div>
	)
};

export default Home;
