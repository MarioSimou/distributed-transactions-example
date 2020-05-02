import React from 'react'
import {
  AppBar,
  Toolbar,
  Typography
} from '@material-ui/core'
import makeStyles from '@material-ui/styles/makeStyles'
import {Link} from 'react-router-dom'
import {useUserProfile} from '../../utils/hooks'
import onClickLogOut from './utils/onClickLogOut.js'

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