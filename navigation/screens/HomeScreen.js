import * as WebBrowser from 'expo-web-browser';
import React from 'react';
import {
  Image,
  Platform,
  ScrollView,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
  Alert, 
  Animated,
} from 'react-native';
import MapView, {Marker} from 'react-native-maps'

import { MonoText } from '../components/StyledText';

export default class HomeScreen extends React.Component {
  
componentDidMount = () => {
    this.getRouteData();

    setInterval(() => {
      this.getTransitVans();
    }, 100);
 }
  
 constructor(props){
    super(props);
    this.state = {
      vans: [],
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
  )};

  getTransitVans = () => {
    fetch('https://xajgj7zaxg.execute-api.us-east-1.amazonaws.com/prod/vans', {
         method: 'GET'
      })
      .then((response) => response.json())
      .then((vans) => {
        if(vans && vans.length > 0){
          vans.forEach(van => {
            van.latitude = Number(van.Latitude);
            van.longitude = Number(van.Longitude);
            van.LatLng = {
              latitude: van.latitude,
              longitude: van.longitude
            }
          });

          this.setState({ vans });
        }
       
      })
      .catch((error) => {
         console.error(error);
      });
  };

  getRouteData = () => {
    fetch('https://xajgj7zaxg.execute-api.us-east-1.amazonaws.com/prod/routes?route_number=1', {
         method: 'GET'
      })
      .then((response) => response.json())
      .then((responseJson) => {
        //console.log(responseJson);

        responseJson.Stops.forEach(stop => {
          stop.LatLng.latitude = Number(stop.LatLng.Latitude);
          stop.LatLng.longitude = Number(stop.LatLng.Longitude);
        });

        const stops = responseJson.Stops;
        
        //console.log(stops);
         this.setState({
          stops,
          people: [],
         })
         
      })
      .catch((error) => {
         console.error(error);
      });
  }

  render() {  
    return (
      <View style={styles.container}>
        
        <MapView style={styles.map} initialRegion={this.state.region} showsUserLocation={true}>
          {this.state.stops.map(marker => (
            <Marker
              key={marker.ID}
              coordinate={marker.LatLng}
              title={marker.Title}
              description={marker.Description}
            />
          ))}
          {this.state.vans.map(marker => (
            <Marker
              key={marker.van_id}
              coordinate={marker.LatLng}
              pinColor='blue'
            />
          ))}
        </MapView>
      </View>
    );
  }
}

HomeScreen.navigationOptions = {
  header: null,
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  contentContainer: {
    paddingTop: 30,
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
  },
  codeHighlightText: {
    color: 'rgba(96,100,109, 0.8)',
  },
  codeHighlightContainer: {
    backgroundColor: 'rgba(0,0,0,0.05)',
    borderRadius: 3,
    paddingHorizontal: 4,
  },
  getStartedText: {
    fontSize: 17,
    color: 'rgba(96,100,109, 1)',
    lineHeight: 24,
    textAlign: 'center',
  },
  tabBarInfoContainer: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    ...Platform.select({
      ios: {
        shadowColor: 'black',
        shadowOffset: { width: 0, height: -3 },
        shadowOpacity: 0.1,
        shadowRadius: 3,
      },
      android: {
        elevation: 20,
      },
    }),
    alignItems: 'center',
    backgroundColor: '#fbfbfb',
    paddingVertical: 20,
  },
  tabBarInfoText: {
    fontSize: 17,
    color: 'rgba(96,100,109, 1)',
    textAlign: 'center',
  },
  navigationFilename: {
    marginTop: 5,
  },
  helpContainer: {
    marginTop: 15,
    alignItems: 'center',
  },
  helpLink: {
    paddingVertical: 15,
  },
  helpLinkText: {
    fontSize: 14,
    color: '#2e78b7',
  },
});
