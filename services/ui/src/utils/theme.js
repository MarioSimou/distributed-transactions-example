import createMuiTheme from '@material-ui/core/styles/createMuiTheme'

const palette = {
  primary: {
    light: '#92c2f2',
    main: '#1976d2',
  },
  // secondary: {
  //   light: '#0066ff',
  //   main: '#0044ff',
  //   contrastText: '#ffcc00',
  // },
}


const theme = createMuiTheme({
  palette,
})

export default theme