using Go = import "capnp/go.capnp";
$Go.package("core");
$Go.import("github.com/fire/go-es_core");

@0x986343ae07c42f59;

struct State {
	union { 
		mouse @0 :Bool;
		kb	@1 :Bool;
		mouseReset @2 :Bool;
		configLookAround @3 :Bool;
	}
	
	quaternion :group {
		w @4 :Float32;
		x @5 :Float32;
		y @6 :Float32;
		z @7 :Float32;
	}
	
	lookAround :group {
		manipulateObject @8 :Bool;
	}
}

struct InputMouse {
	w @0 :Float32;
	x @1 :Float32;
	y @2 :Float32;
	z @3 :Float32;
	buttons @4:UInt32;
}

struct InputKb {
	w @0 :Bool;
	a @1 :Bool;
	s @2 :Bool;
	d @3 :Bool;
	space @4 :Bool;
	lalt @5 :Bool;
}

struct Stop {
	stop @0 :Bool;
}
