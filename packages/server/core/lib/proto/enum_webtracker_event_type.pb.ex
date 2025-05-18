defmodule Realtime.Pb.WebTrackerEventType do
  @moduledoc false
  use Protobuf, enum: true, protoc_gen_elixir_version: "0.14.1", syntax: :proto3

  def descriptor do
    # credo:disable-for-next-line
    %Google.Protobuf.EnumDescriptorProto{
      name: "WebTrackerEventType",
      value: [
        %Google.Protobuf.EnumValueDescriptorProto{
          name: "WEB_TRACKER_EVENT_UNSPECIFIED",
          number: 0,
          options: nil,
          __unknown_fields__: []
        },
        %Google.Protobuf.EnumValueDescriptorProto{
          name: "WEB_TRACKER_PAGE_EXIT",
          number: 1,
          options: nil,
          __unknown_fields__: []
        },
        %Google.Protobuf.EnumValueDescriptorProto{
          name: "WEB_TRACKER_PAGE_VIEW",
          number: 2,
          options: nil,
          __unknown_fields__: []
        },
        %Google.Protobuf.EnumValueDescriptorProto{
          name: "WEB_TRACKER_CLICK",
          number: 3,
          options: nil,
          __unknown_fields__: []
        },
        %Google.Protobuf.EnumValueDescriptorProto{
          name: "WEB_TRACKER_IDENTIFY",
          number: 4,
          options: nil,
          __unknown_fields__: []
        }
      ],
      options: nil,
      reserved_range: [],
      reserved_name: [],
      __unknown_fields__: []
    }
  end

  field :WEB_TRACKER_EVENT_UNSPECIFIED, 0
  field :WEB_TRACKER_PAGE_EXIT, 1
  field :WEB_TRACKER_PAGE_VIEW, 2
  field :WEB_TRACKER_CLICK, 3
  field :WEB_TRACKER_IDENTIFY, 4
end
