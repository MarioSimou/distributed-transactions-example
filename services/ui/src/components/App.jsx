import React from 'react'
import { ThemeProvider } from '@material-ui/core/styles'
import {
  Router,
  Route,
  Switch
} from 'react-router-dom'
import Home from './Home.jsx'
import SignIn from './SignIn.jsx'
import SignUp from './SignUp.jsx'
import CssBaseline from '@material-ui/core/CssBaseline'
import Navbar from './Navbar.jsx'
import theme from '../utils/theme.js'
import history from '../utils/history.js'

const App = () => {

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline>
        <Router history={history}>
          <Navbar/>
          <Switch>
            <Route path="/" exact component={Home} />
            <Route path="/signin" exact component={SignIn} />
            <Route path="/signup" exact component={SignUp} />
          </Switch>
        </Router>
      </CssBaseline>
    </ThemeProvider>
  )
}

export default App