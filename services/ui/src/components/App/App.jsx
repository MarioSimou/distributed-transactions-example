import React from 'react'
import { ThemeProvider } from '@material-ui/core/styles'
import {
  Router,
  Route,
  Switch
} from 'react-router-dom'
import Home from '../Home/Home.jsx'
import SignIn from '../SignIn/SignIn.jsx'
import SignUp from '../SignUp/SignUp.jsx'
import CssBaseline from '@material-ui/core/CssBaseline'
import Navbar from '../shared/Navbar.jsx'
import theme from '../../utils/theme.js'
import history from '../../utils/history.js'
import * as hooks from '../../utils/hooks.js'
import {CancelToken} from 'axios'
import loadUserProfile from './utils/loadUserProfile.js'

const App = () => {
  const [userProfile, setUserProfile] = React.useState(hooks.initUserValues)

  React.useEffect(() => {
    const source = CancelToken.source()
    loadUserProfile({setUserProfile,source})
    return () => source.cancel()
  }, [])

  console.log(userProfile)
  return (
    <ThemeProvider theme={theme}>
      <hooks.UserProfileContext.Provider value={{userProfile, setUserProfile}}>
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
      </hooks.UserProfileContext.Provider>
    </ThemeProvider>
  )
}

export default App