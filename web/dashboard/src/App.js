import React from 'react';

class App extends React.Component {
	constructor(props) {
		super(props);
		this.doms = {};
		this.rem = parseFloat(getComputedStyle(document.documentElement).fontSize);
	}

	componentDidMount() {
		const { doms } = this;
		const rems = parseInt(doms.sidebar.offsetWidth / this.rem) + 1;
		doms.main.style.marginLeft = `${rems}rem`;
	}

	render() {
		return (
			<React.Fragment>
				<header>
					<nav id="sidebar" ref={(e) => { this.doms.sidebar = e; }}>
						<ul className="noul">
							<li><a href="/">Home</a></li>
						</ul>
					</nav>
				</header>
				<main id="main-content" ref={(e) => { this.doms.main = e; }}>
					<p>Main line 1</p>
					<p>Main line 2</p>
					<p>Main line 3</p>
					<p>Main line 4</p>
				</main>
				<footer id="footer" className="txt-right">Padlock TOTP Service &copy; 2019</footer>
			</React.Fragment>
		)
	}
}

export default App;