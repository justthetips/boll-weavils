import React from 'react';
import { StyleSheet, Text, View, Alert, Animated, } from 'react-native';
import MapView, {Marker} from 'react-native-maps'



export default class App extends React.Component {
  getRouteData = () => {
    fetch('https://jsonplaceholder.typicode.com/posts/1', {
         method: 'GET'
      })
      .then((response) => response.json())
      .then((responseJson) => {
          
         this.setState({
          stops: [],
          people: [],
         })
         
      })
      .catch((error) => {
         console.error(error);
      });
  }

  setMyPosition = () => {
    navigator.geolocation.getCurrentPosition(
      (position) => {
        console.log(position);
        this.setState({
          latitude: position.coords.latitude,
          longitude: position.coords.longitude,
          error: null,
        });
      },
      (error) => this.setState({ error: error.message }),
      { enableHighAccuracy: false, timeout: 200000, maximumAge: 1000 },
    );

  }

  componentDidMount = () => {
    this.getRouteData();
 }

  constructor(props){
    super(props);
    this.state = {
      region: {
        latitude: 26.501,
        longitude: -80.232,
        latitudeDelta: .75,
        longitudeDelta: .75,
      },
      stops: [],
      people: [],
      latitude: null,
      longitude: null,
      error: null,
    }
    this.setMyPosition();
  }
  
  onRegionChange(region) {
    this.setState({ region });
  }

  render(){  
    
    return (
      <View style={styles.container}>
        
        <Text>Open up App.js to start working on your app!</Text>
        <MapView style={styles.map} initialRegion={this.state.region} showsUserLocation={true}>
          {this.state.stops.map(marker => (
            <Marker
              key={marker.id}
              coordinate={marker.latlng}
              title={marker.title}
              description={marker.description}
            />
          ))}
          {this.state.people.map(marker => (
            <Marker
              key={marker.id}
              coordinate={marker.latlng}
              title={marker.title}
              description={marker.description}
              pinColor="#000000"
            />
          ))}
        </MapView>
      </View>
    );
  }
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
  },
  markerWrap: {
    alignItems: "center",
    justifyContent: "center",
  },
  marker: {
    width: 8,
    height: 8,
    borderRadius: 4,
    backgroundColor: "rgba(130,4,150, 0.9)",
  },
  ring: {
    width: 24,
    height: 24,
    borderRadius: 12,
    backgroundColor: "rgba(130,4,150, 0.3)",
    position: "absolute",
    borderWidth: 1,
    borderColor: "rgba(130,4,150, 0.5)",
  }, 
  map: {
    ...StyleSheet.absoluteFillObject,
  }
});
