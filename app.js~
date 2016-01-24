var EduApp = angular.module("EduApp", ['ngRoute']);

EduApp.controller('EventDetailsCtrl', [
    '$scope', 
    '$routeParams',
    '$http',
    function($scope, $routeParams, $http){
        $scope.eventId = $routeParams.id;

        $http
            .get('api/locations/' + $scope.eventId)
            .then(function(result){
                $scope.name = result.data.name;
                $scope.venue = result.data.venue;
                $scope.description = result.data.description;
                $scope.date = result.data.date;
                $scope.time = result.data.time;
                $scope.image = result.data.image;
            });
}]);

EduApp.controller('EventsCtrl', [
    '$scope', 
    '$routeParams',
    '$http',
    function($scope, $routeParams, $http){
        
        var latitude = $routeParams.latitude;
        var longitude = $routeParams.longitude;

        $http.get("/api/locations")
                .then(function(response) {
                    $scope.centres = response.data;
                    var centresArray = $scope.centres;
                    var locationMap = {};

                    $scope.centres.forEach(function(centre){
                        locationMap[centre.name] = { latitude : centre.latitude, longitude : centre.longitude, id : centre._id, image: centre.image }
                    });

                    var orderedCentres = geolib.orderByDistance({
                        latitude: latitude,
                        longitude: longitude
                    }, locationMap);

                    $scope.orderedCentres = orderedCentres.map(function(centre){
                        centre.name = centre.key;
                        centre.image = locationMap[centre.name].image;
                        centre.distance = centre.distance / 1000;
                        centre._id = locationMap[centre.name].id;
                        delete centre.key;
                        return centre;
                    })
                });
}]);


EduApp.controller('MainCtrl', ['$scope', '$http', '$routeParams', '$location',
    function($scope, $http, $routeParams, $location) {

        // $scope.center = (geolib.getCenter($scope.venues));
        

        $scope.centres = [];

        $http.get("/api/locations")
            .then(function(response) {
                $scope.centres = response.data;
            });

        // var position;
        $scope.nearMe = function() {
            if (navigator.geolocation) {
                navigator.geolocation.getCurrentPosition($scope.getLocation);
            }
        }

        $scope.findEvents = (function(position) {

            console.log($scope.target_latitude + ", " + $scope.target_longitude);
            var target_latitude = $scope.target_latitude;
            var target_longitude = $scope.target_longitude;
            
            $location.path("/events/" + target_latitude + "/" + target_longitude );
        });

        $scope.getLocation = function(position) {
            //$scope.target_name = position.name;
            $scope.target_latitude = position.coords.latitude;
            $scope.target_longitude = position.coords.longitude;
            $scope.located = true;
            $scope.$apply(); //this triggers a $digest
        };

        $scope.addCentre = function(centre) {

            var theCenter = {
                name: centre.name,
                latitude: centre.latitude,
                longitude: centre.longitude
            };

            $http
                .post('/api/locations', theCenter)
                .then(function(result) {
                    $scope.centres.push(theCenter);
                    $scope.centre = {};
                })
                .catch(function(e) {
                    alert(JSON.stringify(e));
                });
        }

    }
]);

EduApp.config(function($routeProvider){

    $routeProvider.when('/', {
        templateUrl:'templates/home.html',
        controller : 'MainCtrl'
    }).when('/events', {
        templateUrl:'templates/events.html',
        controller : 'MainCtrl'
    })
    .when('/legal', {
        templateUrl:'templates/legal.html',
        controller : 'MainCtrl'
    }).when('/about', {
        templateUrl:'templates/about.html',
        controller : 'MainCtrl'
    }).when('/contact', {
        templateUrl:'templates/contact.html',
        controller : 'MainCtrl'
    }).when('/events/:id', {
        templateUrl:'templates/event_details.html',
        controller : 'EventDetailsCtrl'
    }).when('/events/:latitude/:longitude', {
        templateUrl:'templates/events.html',
        controller : 'EventsCtrl'
    });

});