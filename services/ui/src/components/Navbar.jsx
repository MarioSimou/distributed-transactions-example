import React from 'react'
import {
  AppBar,
  Toolbar,
  Typography
} from '@material-ui/core'
import makeStyles from '@material-ui/styles/makeStyles'
import {Link} from 'react-router-dom'
import {useUserProfile, initUserValues} from '../utils/hooks'
import httpClient from '../utils/httpClient'
import history from '../utils/history'

const onClickLogOut = ({setUserProfile}) => async () => {
  try {
    const uri = new URL(`/api/v1/logout`, process.env.REACT_APP_CUSTOMERS_API)
    const {status} = await httpClient.post(uri.toString())
    if (status === 204){
      setUserProfile(initUserValues)
      history.push('/')
    }
  }catch(e){
    window.alert(e.response && e.response.data && e.response.data.message || e.message)
  }
}

const Login = () => {
  const classes = useStyles()
  const {userProfile, setUserProfile} = useUserProfile()
  const isLoggedIn = Boolean(userProfile.email)

  return (
    <AppBar position="sticky" className={classes.root}>
      <Toolbar className={classes.toolbar} variant="dense">
        <Link to="/" className={classes.link}>
          <Typography variant="h6">Home</Typography>
        </Link>
        {!isLoggedIn &&  
          <Link to="/signin" className={classes.link}>
            <Typography variant="h6">Sign In</Typography>
          </Link>}
        {!isLoggedIn &&
        <Link to="/signup" className={classes.link}>
          <Typography variant="h6">Sign Up</Typography>
        </Link>}
        {isLoggedIn &&
        <Link to="/logout" className={classes.link}>
          <Typography variant="h6" onClick={onClickLogOut({setUserProfile})}>Logout</Typography>
        </Link>}
      </Toolbar>
    </AppBar>
  )
}

const useStyles = makeStyles(theme => ({
  root: {
    width: '100vw',
    minWidth: '100vw',
  },
  toolbar: {
    width: '100%',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
  },
  link: {
    fontWeight: 400,
    color: theme.palette.primary.contrastText,
    marginLeft: theme.spacing(2),
    borderBottom: `1px solid transparent`,
    textDecoration: 'none',
    '&:hover': {      
      borderBottom: `1px solid ${theme.palette.primary.contrastText}`,
    }
  },
}))

export default Login