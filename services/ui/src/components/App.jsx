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
import * as hooks from '../utils/hooks.js'
import httpClient from '../utils/httpClient'
import {CancelToken} from 'axios'

const loadUserProfile = async ({setUserProfile, source}) => {
  try {
    const uri = new URL(`/api/v1/signin`, process.env.REACT_APP_CUSTOMERS_API)
    const {data} = await httpClient.get(uri.toString(), {cancelToken: source.token})
    setUserProfile(data.data)
  }catch(e){
    console.warn('No user profile loaded')
  }
}

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