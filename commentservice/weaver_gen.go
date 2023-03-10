package commentservice

// Code generated by "weaver generate". DO NOT EDIT.
import (
	"context"
	"fmt"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"time"
)

func init() {
	codegen.Register(codegen.Registration{
		Name:     "github.com/edmarfelipe/serviceweaver-example/commentservice/Service",
		Iface:    reflect.TypeOf((*Service)(nil)).Elem(),
		New:      func() any { return &commentService{} },
		ConfigFn: func(i any) any { return i.(*commentService).WithConfig.Config() },
		LocalStubFn: func(impl any, tracer trace.Tracer) any {
			return service_local_stub{impl: impl.(Service), tracer: tracer}
		},
		ClientStubFn: func(stub codegen.Stub, caller string) any {
			return service_client_stub{stub: stub, getByPostMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/edmarfelipe/serviceweaver-example/commentservice/Service", Method: "GetByPost"}), createCommentMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/edmarfelipe/serviceweaver-example/commentservice/Service", Method: "CreateComment"})}
		},
		ServerStubFn: func(impl any, addLoad func(uint64, float64)) codegen.Server {
			return service_server_stub{impl: impl.(Service), addLoad: addLoad}
		},
	})
}

// Local stub implementations.

type service_local_stub struct {
	impl   Service
	tracer trace.Tracer
}

func (s service_local_stub) GetByPost(ctx context.Context, a0 int) (r0 []Comment, err error) {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "commentservice.Service.GetByPost", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.GetByPost(ctx, a0)
}

func (s service_local_stub) CreateComment(ctx context.Context, a0 int, a1 string) (err error) {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "commentservice.Service.CreateComment", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.CreateComment(ctx, a0, a1)
}

// Client stub implementations.

type service_client_stub struct {
	stub                 codegen.Stub
	getByPostMetrics     *codegen.MethodMetrics
	createCommentMetrics *codegen.MethodMetrics
}

func (s service_client_stub) GetByPost(ctx context.Context, a0 int) (r0 []Comment, err error) {
	// Update metrics.
	start := time.Now()
	s.getByPostMetrics.Count.Add(1)

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "commentservice.Service.GetByPost", trace.WithSpanKind(trace.SpanKindClient))
	}

	defer func() {
		// Catch and return any panics detected during encoding/decoding/rpc.
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
		err = s.stub.WrapError(err)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			s.getByPostMetrics.ErrorCount.Add(1)
		}
		span.End()

		s.getByPostMetrics.Latency.Put(float64(time.Since(start).Microseconds()))
	}()

	// Preallocate a buffer of the right size.
	size := 0
	size += 8
	enc := codegen.NewEncoder()
	enc.Reset(size)

	// Encode arguments.
	enc.Int(a0)
	var shardKey uint64

	// Call the remote method.
	s.getByPostMetrics.BytesRequest.Put(float64(len(enc.Data())))
	var results []byte
	results, err = s.stub.Run(ctx, 1, enc.Data(), shardKey)
	if err != nil {
		return
	}
	s.getByPostMetrics.BytesReply.Put(float64(len(results)))

	// Decode the results.
	dec := codegen.NewDecoder(results)
	r0 = serviceweaver_dec_slice_Comment_a6aa7c1a(dec)
	err = dec.Error()
	return
}

func (s service_client_stub) CreateComment(ctx context.Context, a0 int, a1 string) (err error) {
	// Update metrics.
	start := time.Now()
	s.createCommentMetrics.Count.Add(1)

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "commentservice.Service.CreateComment", trace.WithSpanKind(trace.SpanKindClient))
	}

	defer func() {
		// Catch and return any panics detected during encoding/decoding/rpc.
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
		err = s.stub.WrapError(err)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			s.createCommentMetrics.ErrorCount.Add(1)
		}
		span.End()

		s.createCommentMetrics.Latency.Put(float64(time.Since(start).Microseconds()))
	}()

	// Preallocate a buffer of the right size.
	size := 0
	size += 8
	size += (4 + len(a1))
	enc := codegen.NewEncoder()
	enc.Reset(size)

	// Encode arguments.
	enc.Int(a0)
	enc.String(a1)
	var shardKey uint64

	// Call the remote method.
	s.createCommentMetrics.BytesRequest.Put(float64(len(enc.Data())))
	var results []byte
	results, err = s.stub.Run(ctx, 0, enc.Data(), shardKey)
	if err != nil {
		return
	}
	s.createCommentMetrics.BytesReply.Put(float64(len(results)))

	// Decode the results.
	dec := codegen.NewDecoder(results)
	err = dec.Error()
	return
}

// Server stub implementations.

type service_server_stub struct {
	impl    Service
	addLoad func(key uint64, load float64)
}

// GetStubFn implements the stub.Server interface.
func (s service_server_stub) GetStubFn(method string) func(ctx context.Context, args []byte) ([]byte, error) {
	switch method {
	case "GetByPost":
		return s.getByPost
	case "CreateComment":
		return s.createComment
	default:
		return nil
	}
}

func (s service_server_stub) getByPost(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 int
	a0 = dec.Int()

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, appErr := s.impl.GetByPost(ctx, a0)

	// Encode the results.
	enc := codegen.NewEncoder()
	serviceweaver_enc_slice_Comment_a6aa7c1a(enc, r0)
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s service_server_stub) createComment(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 int
	a0 = dec.Int()
	var a1 string
	a1 = dec.String()

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	appErr := s.impl.CreateComment(ctx, a0, a1)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.Error(appErr)
	return enc.Data(), nil
}

// AutoMarshal implementations.

var _ codegen.AutoMarshal = &Comment{}

func (x *Comment) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Comment.WeaverMarshal: nil receiver"))
	}
	enc.Int(x.ID)
	enc.Int(x.PostID)
	enc.String(x.Content)
	enc.EncodeBinaryMarshaler(&x.CreateAt)
}

func (x *Comment) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Comment.WeaverUnmarshal: nil receiver"))
	}
	x.ID = dec.Int()
	x.PostID = dec.Int()
	x.Content = dec.String()
	dec.DecodeBinaryUnmarshaler(&x.CreateAt)
}

// Encoding/decoding implementations.

func serviceweaver_enc_slice_Comment_a6aa7c1a(enc *codegen.Encoder, arg []Comment) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		(arg[i]).WeaverMarshal(enc)
	}
}

func serviceweaver_dec_slice_Comment_a6aa7c1a(dec *codegen.Decoder) []Comment {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]Comment, n)
	for i := 0; i < n; i++ {
		(&res[i]).WeaverUnmarshal(dec)
	}
	return res
}
