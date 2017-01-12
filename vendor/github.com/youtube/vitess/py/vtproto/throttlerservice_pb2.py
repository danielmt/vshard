# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: throttlerservice.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


import throttlerdata_pb2 as throttlerdata__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='throttlerservice.proto',
  package='throttlerservice',
  syntax='proto3',
  serialized_pb=_b('\n\x16throttlerservice.proto\x12\x10throttlerservice\x1a\x13throttlerdata.proto2\xaf\x01\n\tThrottler\x12M\n\x08MaxRates\x12\x1e.throttlerdata.MaxRatesRequest\x1a\x1f.throttlerdata.MaxRatesResponse\"\x00\x12S\n\nSetMaxRate\x12 .throttlerdata.SetMaxRateRequest\x1a!.throttlerdata.SetMaxRateResponse\"\x00\x62\x06proto3')
  ,
  dependencies=[throttlerdata__pb2.DESCRIPTOR,])
_sym_db.RegisterFileDescriptor(DESCRIPTOR)





import abc
from grpc.beta import implementations as beta_implementations
from grpc.framework.common import cardinality
from grpc.framework.interfaces.face import utilities as face_utilities

class BetaThrottlerServicer(object):
  """<fill me in later!>"""
  __metaclass__ = abc.ABCMeta
  @abc.abstractmethod
  def MaxRates(self, request, context):
    raise NotImplementedError()
  @abc.abstractmethod
  def SetMaxRate(self, request, context):
    raise NotImplementedError()

class BetaThrottlerStub(object):
  """The interface to which stubs will conform."""
  __metaclass__ = abc.ABCMeta
  @abc.abstractmethod
  def MaxRates(self, request, timeout):
    raise NotImplementedError()
  MaxRates.future = None
  @abc.abstractmethod
  def SetMaxRate(self, request, timeout):
    raise NotImplementedError()
  SetMaxRate.future = None

def beta_create_Throttler_server(servicer, pool=None, pool_size=None, default_timeout=None, maximum_timeout=None):
  import throttlerdata_pb2
  import throttlerdata_pb2
  import throttlerdata_pb2
  import throttlerdata_pb2
  request_deserializers = {
    ('throttlerservice.Throttler', 'MaxRates'): throttlerdata_pb2.MaxRatesRequest.FromString,
    ('throttlerservice.Throttler', 'SetMaxRate'): throttlerdata_pb2.SetMaxRateRequest.FromString,
  }
  response_serializers = {
    ('throttlerservice.Throttler', 'MaxRates'): throttlerdata_pb2.MaxRatesResponse.SerializeToString,
    ('throttlerservice.Throttler', 'SetMaxRate'): throttlerdata_pb2.SetMaxRateResponse.SerializeToString,
  }
  method_implementations = {
    ('throttlerservice.Throttler', 'MaxRates'): face_utilities.unary_unary_inline(servicer.MaxRates),
    ('throttlerservice.Throttler', 'SetMaxRate'): face_utilities.unary_unary_inline(servicer.SetMaxRate),
  }
  server_options = beta_implementations.server_options(request_deserializers=request_deserializers, response_serializers=response_serializers, thread_pool=pool, thread_pool_size=pool_size, default_timeout=default_timeout, maximum_timeout=maximum_timeout)
  return beta_implementations.server(method_implementations, options=server_options)

def beta_create_Throttler_stub(channel, host=None, metadata_transformer=None, pool=None, pool_size=None):
  import throttlerdata_pb2
  import throttlerdata_pb2
  import throttlerdata_pb2
  import throttlerdata_pb2
  request_serializers = {
    ('throttlerservice.Throttler', 'MaxRates'): throttlerdata_pb2.MaxRatesRequest.SerializeToString,
    ('throttlerservice.Throttler', 'SetMaxRate'): throttlerdata_pb2.SetMaxRateRequest.SerializeToString,
  }
  response_deserializers = {
    ('throttlerservice.Throttler', 'MaxRates'): throttlerdata_pb2.MaxRatesResponse.FromString,
    ('throttlerservice.Throttler', 'SetMaxRate'): throttlerdata_pb2.SetMaxRateResponse.FromString,
  }
  cardinalities = {
    'MaxRates': cardinality.Cardinality.UNARY_UNARY,
    'SetMaxRate': cardinality.Cardinality.UNARY_UNARY,
  }
  stub_options = beta_implementations.stub_options(host=host, metadata_transformer=metadata_transformer, request_serializers=request_serializers, response_deserializers=response_deserializers, thread_pool=pool, thread_pool_size=pool_size)
  return beta_implementations.dynamic_stub(channel, 'throttlerservice.Throttler', cardinalities, options=stub_options)
# @@protoc_insertion_point(module_scope)
